from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[1]

class ComposeProductionSecurityTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.compose = (ROOT / "compose.yaml").read_text(encoding="utf-8")
        cls.workflow = (ROOT / ".github" / "workflows" / "web-stack.yml").read_text(encoding="utf-8")

    def test_game_api_host_port_is_loopback_only(self):
        expected = '"127.0.0.1:' + chr(36) + '{GAME_API_PORT:-18080}:8080"'
        self.assertIn(expected, self.compose)
        self.assertNotIn('      - "' + chr(36) + '{GAME_API_PORT:-18080}:8080"', self.compose)

    def test_database_passwords_are_required_not_defaulted(self):
        pg_required = chr(36) + '{POSTGRES_PASSWORD:?POSTGRES_PASSWORD must be set}'
        exporter_required = chr(36) + '{POSTGRES_EXPORTER_PASSWORD:?POSTGRES_EXPORTER_PASSWORD must be set}'
        self.assertIn(pg_required, self.compose)
        self.assertIn(exporter_required, self.compose)
        self.assertNotIn(chr(36) + '{POSTGRES_PASSWORD:-zneondrive-local}', self.compose)
        self.assertNotIn(chr(36) + '{POSTGRES_EXPORTER_PASSWORD:-zneondrive-monitor-local}', self.compose)

    def test_web_stack_ci_supplies_test_credentials(self):
        self.assertIn('POSTGRES_PASSWORD: ci-postgres-test-password-not-for-production', self.workflow)
        self.assertIn('POSTGRES_EXPORTER_PASSWORD: ci-exporter-test-password-not-for-production', self.workflow)

if __name__ == "__main__":
    unittest.main()
