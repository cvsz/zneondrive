"""Regression checks for the public AdSense integration boundaries."""
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[1]
PUBLISHER = "ca-pub-4971034675329740"
LOADER = "https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"


class AdSenseSurfaceTests(unittest.TestCase):
    def test_only_editorial_homepage_loads_adsense(self):
        homepage = (ROOT / "web/index.html").read_text(encoding="utf-8")
        head = homepage.split("<head>", 1)[1].split("</head>", 1)[0]
        self.assertEqual(head.count(PUBLISHER), 1)
        self.assertEqual(head.count(LOADER), 1)
        for path in ("web/privacy.html", "web/terms.html", "web/client-v-0-0-1/index.html"):
            with self.subTest(path=path):
                content = (ROOT / path).read_text(encoding="utf-8")
                self.assertNotIn(LOADER, content)
                self.assertNotIn(PUBLISHER, content)

    def test_privacy_notice_discloses_google_advertising(self):
        privacy = (ROOT / "web/privacy.html").read_text(encoding="utf-8")
        for term in ("Google AdSense", "cookies", "consent", "policies.google.com/technologies/ads"):
            with self.subTest(term=term):
                self.assertIn(term, privacy)

    def test_nginx_csp_permits_adsense_loader(self):
        nginx = (ROOT / "web/nginx.conf").read_text(encoding="utf-8")
        csp = next(line for line in nginx.splitlines() if "Content-Security-Policy" in line)
        self.assertIn("script-src 'self' https://pagead2.googlesyndication.com", csp)
        self.assertIn("frame-src https://*.googlesyndication.com", csp)
        self.assertIn("frame-ancestors 'none'", csp)
        self.assertIn("location ^~ /api/", nginx)


if __name__ == "__main__":
    unittest.main()
