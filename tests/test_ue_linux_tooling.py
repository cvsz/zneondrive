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
        self.assertIn("UE_CLIENT_TARGET=present", result.stdout)
        self.assertIn("UE_SERVER_TARGET=present", result.stdout)
        self.assertRegex(result.stdout, r"UE_FREE_DISK_BYTES=\d+")

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
