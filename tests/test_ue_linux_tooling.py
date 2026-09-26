import json
import os
import subprocess
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "tools" / "ue-linux.sh"
PROJECT = ROOT / "game" / "NeonDrive.uproject"


class UELinuxToolingTests(unittest.TestCase):
    def setUp(self) -> None:
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.engine = Path(self.tmp.name) / "UE"
        (self.engine / "Engine/Build/BatchFiles/Linux").mkdir(parents=True)
        (self.engine / "Engine/Binaries/DotNET/UnrealBuildTool").mkdir(parents=True)
        (self.engine / "Engine/Binaries/ThirdParty/DotNet/10.0.0/linux").mkdir(parents=True)
        (self.engine / "Engine/Build/Build.version").write_text(
            json.dumps({"MajorVersion": 5, "MinorVersion": 8, "PatchVersion": 2}),
            encoding="utf-8",
        )
        self.calls = Path(self.tmp.name) / "calls.log"
        self._write_exe(
            self.engine / "Engine/Build/BatchFiles/Linux/Build.sh",
            '#!/usr/bin/env bash\necho "build:$*" >> "$CALLS_LOG"\n',
        )
        self._write_exe(
            self.engine / "Engine/Build/BatchFiles/RunUAT.sh",
            '#!/usr/bin/env bash\necho "uat:$*" >> "$CALLS_LOG"\n',
        )

    def _write_exe(self, path: Path, content: str) -> None:
        path.write_text(content, encoding="utf-8")
        path.chmod(0o755)

    def _write_installed_manifests(self, target: str, modules=("Core", "CoreUObject", "TraceLog")) -> None:
        self._write_exe(self.engine / "Engine/Build/InstalledBuild.txt", "")
        ubt = self.engine / "Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll"
        ubt.write_text("fixture", encoding="utf-8")
        for module in modules:
            path = self.engine / f"Engine/Intermediate/Build/Linux/x64/{target}/Development/{module}/{module}.precompiled"
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text("fixture", encoding="utf-8")

    def _run(self, *args: str, check: bool = True, extra_env=None) -> subprocess.CompletedProcess[str]:
        env = os.environ.copy()
        env.update(
            UE_ROOT=str(self.engine),
            ZNEON_UE_PROJECT=str(PROJECT),
            CALLS_LOG=str(self.calls),
            UE_MIN_FREE_GB="0",
        )
        if extra_env:
            env.update(extra_env)
        return subprocess.run(
            ["bash", str(SCRIPT), *args],
            env=env,
            text=True,
            capture_output=True,
            check=check,
        )

    def test_source_tree_generate_script_is_preferred(self) -> None:
        self._write_exe(
            self.engine / "GenerateProjectFiles.sh",
            '#!/usr/bin/env bash\necho "gpf:$*" >> "$CALLS_LOG"\n',
        )
        self._run("generate")
        text = self.calls.read_text(encoding="utf-8")
        self.assertIn("gpf:-project=", text)
        self.assertIn("NeonDrive.uproject", text)
        self.assertIn(" -game", text)

    def test_source_tree_preflight_reports_build_readiness(self) -> None:
        self._write_exe(
            self.engine / "GenerateProjectFiles.sh",
            '#!/usr/bin/env bash\nexit 0\n',
        )
        result = self._run("preflight")
        self.assertIn("PREFLIGHT_STATUS=ok", result.stdout)
        self.assertIn("UE_VERSION=5.8.2", result.stdout)
        self.assertIn("UE_PROJECTFILES_MODE=source-script", result.stdout)
        self.assertIn("UE_CLIENT_SERVER_PREFLIGHT=source-tree-compile", result.stdout)
        self.assertIn("UE_CLIENT_TARGET=present", result.stdout)
        self.assertIn("UE_SERVER_TARGET=present", result.stdout)
        self.assertRegex(result.stdout, r"UE_FREE_DISK_BYTES=\d+")

    def test_editor_only_installed_build_fails_client_server_preflight(self) -> None:
        editor_manifest = (
            self.engine
            / "Engine/Intermediate/Build/Linux/x64/UnrealEditor/Development/Core/Core.precompiled"
        )
        editor_manifest.parent.mkdir(parents=True, exist_ok=True)
        editor_manifest.write_text("fixture", encoding="utf-8")
        self._write_exe(self.engine / "Engine/Build/InstalledBuild.txt", "")
        (self.engine / "Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll").write_text(
            "fixture", encoding="utf-8"
        )

        result = self._run("preflight", check=False)

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Installed UE build lacks Linux UnrealClient Development precompiled manifests", result.stderr)
        self.assertIn("InstalledBuild.txt alone does not establish Client/Server support", result.stderr)
        self.assertNotIn("UnrealClient Linux Development", self.calls.read_text(encoding="utf-8") if self.calls.exists() else "")

    def test_client_server_capable_installed_build_passes_preflight(self) -> None:
        self._write_installed_manifests("UnrealClient")
        self._write_installed_manifests("UnrealServer")

        result = self._run("preflight")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("UE_CLIENT_SERVER_PREFLIGHT=installed-manifest-baseline-present", result.stdout)

    def test_incomplete_installed_build_reports_missing_required_modules(self) -> None:
        self._write_installed_manifests("UnrealClient", modules=("Core",))

        result = self._run("preflight", check=False)

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Missing: CoreUObject TraceLog", result.stderr)
        self.assertIn("UnrealClient/Development/<Module>/<Module>.precompiled", result.stderr)

    def test_installed_build_uses_bundled_dotnet_and_ubt(self) -> None:
        dll = self.engine / "Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll"
        dll.write_text("fixture", encoding="utf-8")
        self._write_exe(
            self.engine / "Engine/Binaries/ThirdParty/DotNet/10.0.0/linux/dotnet",
            '#!/usr/bin/env bash\necho "dotnet:$*" >> "$CALLS_LOG"\n',
        )
        result = self._run("info")
        self.assertIn("UE_VERSION=5.8.2", result.stdout)
        self.assertIn("UE_PROJECTFILES_MODE=installed-ubt", result.stdout)
        preflight = self._run("preflight")
        self.assertIn("PREFLIGHT_STATUS=ok", preflight.stdout)
        self.assertIn("UE_PROJECTFILES_MODE=installed-ubt", preflight.stdout)
        self.assertIn("UE_CLIENT_SERVER_PREFLIGHT=unverified-ubt-layout", preflight.stdout)
        self._run("generate")
        text = self.calls.read_text(encoding="utf-8")
        self.assertIn("UnrealBuildTool.dll -projectfiles", text)
        self.assertIn("-game -engine", text)

    def test_preflight_rejects_invalid_disk_threshold(self) -> None:
        self._write_exe(
            self.engine / "GenerateProjectFiles.sh",
            '#!/usr/bin/env bash\nexit 0\n',
        )
        result = self._run("preflight", check=False, extra_env={"UE_MIN_FREE_GB": "many"})
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("UE_MIN_FREE_GB must be a non-negative integer", result.stderr)

    def test_wrong_engine_minor_is_rejected(self) -> None:
        (self.engine / "Engine/Build/Build.version").write_text(
            json.dumps({"MajorVersion": 5, "MinorVersion": 7, "PatchVersion": 6}),
            encoding="utf-8",
        )
        result = self._run("info", check=False)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Unreal Engine 5.8.x required", result.stderr)

    def test_build_target_delegates_to_linux_build_script(self) -> None:
        self._run("build-target", "NeonDriveClient")
        text = self.calls.read_text(encoding="utf-8")
        self.assertIn("build:NeonDriveClient Linux Development", text)
        self.assertIn("NeonDrive.uproject", text)
        self.assertIn("-WaitMutex", text)

    def test_build_target_rejects_missing_installed_manifests_before_ubt(self) -> None:
        self._write_exe(self.engine / "Engine/Build/InstalledBuild.txt", "")

        result = self._run("build-target", "NeonDriveClient", check=False)

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Installed UE build lacks Linux UnrealClient Development precompiled manifests", result.stderr)
        self.assertFalse(self.calls.exists())

    def test_package_client_preserves_existing_output_when_target_is_unsupported(self) -> None:
        self._write_exe(self.engine / "Engine/Build/InstalledBuild.txt", "")
        dist = Path(self.tmp.name) / "dist"
        output = dist / "packages/client-linux"
        output.mkdir(parents=True)
        marker = output / "preserve.txt"
        marker.write_text("user data", encoding="utf-8")
        env = {
            "UE_ROOT": str(self.engine),
            "ZNEON_UE_PROJECT": str(PROJECT),
            "DIST_DIR": str(dist),
            "CALLS_LOG": str(self.calls),
        }

        result = subprocess.run(
            ["bash", str(SCRIPT), "package-client"],
            env={**os.environ, **env},
            text=True,
            capture_output=True,
            check=False,
        )

        self.assertNotEqual(result.returncode, 0)
        self.assertTrue(marker.is_file())
        self.assertFalse(self.calls.exists())

    def test_incomplete_installation_fails_required_tool_preflight(self) -> None:
        (self.engine / "Engine/Build/BatchFiles/RunUAT.sh").unlink()

        result = self._run("preflight", check=False)

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("RunUAT.sh missing", result.stderr)

    def test_packaging_keeps_client_and_server_outputs_separate(self) -> None:
        dist = Path(self.tmp.name) / "dist"
        env = os.environ.copy()
        env.update(
            UE_ROOT=str(self.engine),
            ZNEON_UE_PROJECT=str(PROJECT),
            DIST_DIR=str(dist),
            CALLS_LOG=str(self.calls),
        )
        subprocess.run(
            ["bash", str(SCRIPT), "package-all"],
            env=env,
            text=True,
            capture_output=True,
            check=True,
        )
        text = self.calls.read_text(encoding="utf-8")
        self.assertIn(str(dist / "packages/client-linux"), text)
        self.assertIn(str(dist / "packages/server-linux"), text)
        self.assertIn("-client", text)
        self.assertIn("-server -noclient", text)


if __name__ == "__main__":
    unittest.main()
