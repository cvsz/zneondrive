#!/usr/bin/env python3
"""Generate the PROJECT: NEON DRIVE launch character set.

Outputs (repository-relative):
  web/assets/characters/<id>.svg   12 stylized neon portraits (256x256)
  web/characters.json              machine-readable roster for the site gallery

Run:  python3 tools/generate_characters.py
"""
from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
OUT_DIR = ROOT / "web" / "assets" / "characters"
JSON_OUT = ROOT / "web" / "characters.json"

ROSTER = [
    dict(id="kai", th="ไค", en="KAI", species="human", kind="ชาย · นักแข่งรุ่นเก๋า",
         role="Veteran Racer / Garage Chief", faction="Free Roads",
         accent="#22d3ee", spd=8, hdl=7, arm=5, tec=6, chr=8,
         th_bio="อดีตแชมป์สนามใต้ดิน ผันตัวมาเปิดอู่และปั้นนักแข่งรุ่นใหม่",
         en_bio="Ex-underground champion turned garage chief mentoring rookie runners."),
    dict(id="mira", th="มิร่า", en="MIRA", species="human", kind="หญิง · วิศวกร",
         role="Race Engineer / Navigator", faction="Mechanist Guild",
         accent="#f0f", spd=6, hdl=8, arm=4, tec=10, chr=7,
         th_bio="วิศวกรจูนรถอัจฉริยะ อ่านสนามได้ก่อนใครผ่านข้อมูลเทเลเมทรี",
         en_bio="Genius tune engineer who reads the track through live telemetry."),
    dict(id="po", th="โป้", en="PO", species="human", kind="เด็กโต · รุกกี้",
         role="Rookie Runner", faction="Free Roads",
         accent="#a3e635", spd=7, hdl=9, arm=3, tec=5, chr=9,
         th_bio="เด็กส่งของที่ขับแซงทุกคนในซอย ฝันอยากเป็นแชมป์ NOVA CITY",
         en_bio="Courier kid who outruns everyone in the alleys; future champion."),
    dict(id="lungvolt", th="ลุงโวลต์", en="UNCLE VOLT", species="human", kind="คนแก่ · ตำนานช่าง",
         role="Retired Legend / Master Mechanic", faction="Mechanist Guild",
         accent="#fbbf24", spd=4, hdl=6, arm=6, tec=9, chr=10,
         th_bio="ตำนานสนามยุคบุกเบิก ซ่อมได้ทุกอย่างด้วยมือและประสบการณ์",
         en_bio="Pioneer-era legend who can fix anything with hands and memory."),
    dict(id="jump", th="จั๊มพ์", en="JUMP", species="rabbit", kind="กระต่าย · หน่วยสอดแนม",
         role="Sprinter / Scout", faction="Iron Wolves",
         accent="#fb7185", spd=10, hdl=9, arm=3, tec=5, chr=7,
         th_bio="กระต่ายสายฟ้า เร็วที่สุดในตรอกแคบ ไม่มีใครไล่ทัน",
         en_bio="Lightning rabbit; unbeatable in tight alley sprints."),
    dict(id="shell", th="เชลล์", en="SHELL", species="turtle", kind="เต่า · รถถัง",
         role="Heavy Defender / Tanker", faction="NOVA Authority",
         accent="#34d399", spd=3, hdl=5, arm=10, tec=6, chr=5,
         th_bio="กำแพงเคลื่อนที่ของหน่วยงาน เกราะหนา ใจเย็น ไม่เคยถอย",
         en_bio="The Authority's moving wall; patient, armored, immovable."),
    dict(id="sky", th="สกาย", en="SKY", species="bird", kind="นก · ลาดตระเวน",
         role="Aerial Recon", faction="VANTEX",
         accent="#38bdf8", spd=9, hdl=8, arm=3, tec=8, chr=6,
         th_bio="นกนางแอ่นนักสืบอากาศ เห็นทุกความเคลื่อนไหวจากด้านบน",
         en_bio="Swallow recon operative watching every move from above."),
    dict(id="rooster", th="รูสเตอร์", en="ROOSTER", species="chicken", kind="ไก่ · โชว์แมน",
         role="Showman / Trick Driver", faction="VANTEX",
         accent="#f97316", spd=7, hdl=10, arm=4, tec=5, chr=10,
         th_bio="ขวัญใจสนาม ชอบโชว์ลีลาเสี่ยงตายเรียกเสียงเชียร์",
         en_bio="Crowd-favorite daredevil showman of the circuit."),
    dict(id="neoncat", th="นีออน", en="NEON", species="cat", kind="แมว · สายลับ",
         role="Agile Infiltrator", faction="Iron Wolves",
         accent="#a78bfa", spd=9, hdl=10, arm=3, tec=7, chr=8,
         th_bio="แมวลอบเร้น เงียบ ไว และหายตัวได้ก่อนใครจะรู้ตัว",
         en_bio="Silent feline infiltrator; gone before anyone notices."),
    dict(id="tusk", th="ทัสก์", en="TUSK", species="elephant", kind="ช้าง · ขนส่งหนัก",
         role="Heavy Hauler", faction="NOVA Authority",
         accent="#94a3b8", spd=4, hdl=5, arm=10, tec=6, chr=7,
         th_bio="ยักษ์ใหญ่ใจดี แบกของหนักและปกป้องขบวนได้ทั้งทีม",
         en_bio="Gentle giant hauling heavy loads and shielding convoys."),
    dict(id="dash", th="แดช", en="DASH", species="horse", kind="ม้า · นักแข่งมาราธอน",
         role="Endurance Racer", faction="Free Roads",
         accent="#ef4444", spd=9, hdl=7, arm=5, tec=5, chr=8,
         th_bio="ม้าศึกทางไกล ยิ่งวิ่งยิ่งแรง ไม่มีคำว่าหมดในพจนานุกรม",
         en_bio="Long-distance warhorse; stronger with every lap."),
    dict(id="venom", th="เวน่อม", en="VENOM", species="snake", kind="งู · ก่อวินาศกรรม",
         role="Saboteur", faction="Iron Wolves",
         accent="#4ade80", spd=6, hdl=8, arm=4, tec=9, chr=4,
         th_bio="งูพิษแห่งเงามืด วางกับดักและทำลายแผนคู่แข่งอย่างเงียบเชียบ",
         en_bio="Shadow serpent sabotaging rival plans without a trace."),
]


def svg_open(accent: str) -> list[str]:
    return [
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256">',
        "<defs>",
        '<linearGradient id="bg" x1="0" y1="0" x2="0" y2="1">'
        '<stop offset="0" stop-color="#0b1026"/><stop offset="1" stop-color="#05070f"/>'
        "</linearGradient>",
        f'<filter id="glow" x="-40%" y="-40%" width="180%" height="180%">'
        '<feGaussianBlur stdDeviation="3.2" result="b"/>'
        '<feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>'
        "</filter>",
        "</defs>",
        f'<rect x="4" y="4" width="248" height="248" rx="28" fill="url(#bg)" stroke="{accent}" stroke-width="5"/>',
        f'<ellipse cx="128" cy="222" rx="70" ry="14" fill="{accent}" opacity="0.22" filter="url(#glow)"/>',
    ]


def svg_body(accent: str) -> list[str]:
    return [
        # shoulders / jacket
        '<path d="M52 256 Q56 196 96 188 L160 188 Q200 196 204 256 Z" fill="#141a33"/>',
        f'<path d="M96 188 L112 256 M160 188 L144 256" stroke="{accent}" stroke-width="4" opacity="0.7"/>',
        '<path d="M104 196 h48 l-6 60 h-36 Z" fill="#1f2745"/>',
        f'<circle cx="128" cy="216" r="5" fill="{accent}" filter="url(#glow)"/>',
        # scarf
        f'<path d="M84 178 Q128 196 172 178 L168 192 Q128 208 88 192 Z" fill="{accent}" opacity="0.85"/>',
    ]


def svg_visor(accent: str, y: int = 118, w: int = 76) -> list[str]:
    x = 128 - w // 2
    return [
        f'<rect x="{x}" y="{y}" width="{w}" height="22" rx="11" fill="#04060d" stroke="{accent}" stroke-width="3"/>',
        f'<rect x="{x + 8}" y="{y + 6}" width="{w - 16}" height="10" rx="5" fill="{accent}" filter="url(#glow)"/>',
    ]


def svg_name(en: str, th: str, accent: str) -> list[str]:
    return [
        f'<text x="128" y="52" text-anchor="middle" font-family="sans-serif" font-size="21" font-weight="bold" fill="#f4f6ff" letter-spacing="3">{en}</text>',
        f'<text x="128" y="70" text-anchor="middle" font-family="sans-serif" font-size="13" fill="{accent}">{th}</text>',
    ]


def head(shape: str, skin: str, accent: str) -> list[str]:
    parts: list[str] = []
    if shape == "human":
        parts.append('<ellipse cx="128" cy="128" rx="44" ry="48" fill="' + skin + '"/>')
    elif shape == "long":  # horse
        parts.append('<ellipse cx="128" cy="132" rx="38" ry="52" fill="' + skin + '"/>')
        parts.append('<ellipse cx="128" cy="168" rx="20" ry="18" fill="' + skin + '" stroke="#00000033"/>')
    elif shape == "round":  # turtle / elephant base
        parts.append('<circle cx="128" cy="130" r="46" fill="' + skin + '"/>')
    elif shape == "triangle":  # cat
        parts.append('<path d="M84 128 Q88 84 128 82 Q168 84 172 128 Q170 168 128 170 Q86 168 84 128 Z" fill="' + skin + '"/>')
    elif shape == "slim":  # bird / chicken / rabbit / snake
        parts.append('<ellipse cx="128" cy="130" rx="40" ry="46" fill="' + skin + '"/>')
    return parts


def features(sp: str, c: dict) -> list[str]:
    a = c["accent"]
    p: list[str] = []
    if sp == "human":
        if c["id"] == "kai":
            p += ['<path d="M84 118 Q86 78 128 76 Q170 78 172 118 L166 112 Q160 88 128 88 Q96 88 90 112 Z" fill="#1a1d2e"/>',
                  '<path d="M90 112 Q128 96 166 112" stroke="#22d3ee" stroke-width="4" fill="none" opacity="0.8"/>']
            # stubble
            p += ['<path d="M100 150 Q128 162 156 150 L152 162 Q128 170 104 162 Z" fill="#00000022"/>']
        elif c["id"] == "mira":
            p += ['<path d="M82 128 Q80 76 128 74 Q176 76 174 128 L164 128 Q164 92 128 90 Q92 92 92 128 Z" fill="#232742"/>',
                  '<path d="M96 84 Q128 74 160 84" stroke="#f0f" stroke-width="5" fill="none" opacity="0.9"/>',
                  f'<circle cx="170" cy="140" r="6" fill="{a}" filter="url(#glow)"/>']
        elif c["id"] == "po":
            p += ['<path d="M84 112 Q88 70 128 70 Q168 70 172 112 L172 100 L84 100 Z" fill="#a3e635"/>',
                  '<rect x="84" y="100" width="88" height="10" rx="5" fill="#65a30d"/>',
                  '<rect x="150" y="104" width="30" height="6" rx="3" fill="#365314"/>']
        elif c["id"] == "lungvolt":
            p += ['<path d="M92 96 Q128 82 164 96 L160 88 Q128 74 96 88 Z" fill="#9ca3af"/>',
                  '<path d="M98 148 Q128 158 158 148 L154 176 Q128 186 102 176 Z" fill="#d1d5db" opacity="0.9"/>',
                  '<path d="M104 108 h14 M138 108 h14" stroke="#6b7280" stroke-width="3"/>']
    elif sp == "rabbit":
        p += [f'<ellipse cx="102" cy="52" rx="12" ry="34" fill="#e8e4da" stroke="{a}" stroke-width="3"/>',
              f'<ellipse cx="154" cy="52" rx="12" ry="34" fill="#e8e4da" stroke="{a}" stroke-width="3"/>',
              '<ellipse cx="102" cy="54" rx="5" ry="22" fill="#f9a8d4"/>',
              '<ellipse cx="154" cy="54" rx="5" ry="22" fill="#f9a8d4"/>',
              '<ellipse cx="128" cy="150" rx="7" ry="5" fill="#f9a8d4"/>']
    elif sp == "turtle":
        p += [f'<path d="M60 130 Q64 60 128 58 Q192 60 196 130 Q160 108 128 108 Q96 108 60 130 Z" fill="#0f766e" stroke="{a}" stroke-width="4"/>',
              '<circle cx="104" cy="86" r="7" fill="#134e4a"/><circle cx="128" cy="78" r="7" fill="#134e4a"/><circle cx="152" cy="86" r="7" fill="#134e4a"/>',
              '<circle cx="104" cy="146" r="4" fill="#0f766e"/><circle cx="152" cy="146" r="4" fill="#0f766e"/>']
    elif sp == "bird":
        p += [f'<path d="M108 84 Q128 56 148 84 Q138 78 128 86 Q118 78 108 84 Z" fill="{a}"/>',
              '<path d="M118 88 Q128 74 138 88" stroke="#0ea5e9" stroke-width="4" fill="none"/>',
              '<path d="M118 146 L138 146 L128 160 Z" fill="#fbbf24"/>']
    elif sp == "chicken":
        p += ['<circle cx="112" cy="82" r="9" fill="#ef4444"/><circle cx="128" cy="76" r="10" fill="#ef4444"/><circle cx="144" cy="82" r="9" fill="#ef4444"/>',
              '<path d="M116 146 L140 146 L128 162 Z" fill="#fbbf24"/>',
              '<ellipse cx="128" cy="168" rx="7" ry="10" fill="#ef4444"/>']
    elif sp == "cat":
        p += [f'<path d="M88 100 L78 58 L116 84 Z" fill="#2b3352" stroke="{a}" stroke-width="3"/>',
              f'<path d="M168 100 L178 58 L140 84 Z" fill="#2b3352" stroke="{a}" stroke-width="3"/>',
              '<path d="M92 90 L86 68 L108 82 Z" fill="#f9a8d4"/>',
              '<path d="M164 90 L170 68 L148 82 Z" fill="#f9a8d4"/>',
              '<path d="M96 142 h20 M140 142 h20 M98 150 h18 M142 150 h18" stroke="#e5e7eb" stroke-width="2" opacity="0.8"/>',
              '<path d="M124 146 h8 l-4 6 Z" fill="#f9a8d4"/>']
    elif sp == "elephant":
        p += ['<ellipse cx="80" cy="120" rx="16" ry="26" fill="#8fa1b8"/>',
              '<ellipse cx="176" cy="120" rx="16" ry="26" fill="#8fa1b8"/>',
              '<path d="M118 140 Q116 176 128 184 Q140 176 138 140 L132 140 Q132 168 128 172 Q124 168 124 140 Z" fill="#8fa1b8"/>',
              '<path d="M96 108 Q104 102 112 108 M144 108 Q152 102 160 108" stroke="#475569" stroke-width="3" fill="none"/>']
    elif sp == "horse":
        p += [f'<path d="M96 84 L86 44 L118 70 Z" fill="#3f2d23" stroke="{a}" stroke-width="3"/>',
              f'<path d="M160 84 L170 44 L138 70 Z" fill="#3f2d23" stroke="{a}" stroke-width="3"/>',
              '<path d="M128 78 Q150 96 148 140" stroke="#1f2937" stroke-width="12" fill="none"/>',
              '<ellipse cx="120" cy="168" rx="5" ry="7" fill="#0f172a"/><ellipse cx="136" cy="168" rx="5" ry="7" fill="#0f172a"/>']
    elif sp == "snake":
        p += [f'<path d="M78 130 Q70 70 128 64 Q186 70 178 130 Q160 100 128 100 Q96 100 78 130 Z" fill="#14532d" stroke="{a}" stroke-width="4"/>',
              '<circle cx="104" cy="88" r="5" fill="#052e16"/><circle cx="128" cy="82" r="5" fill="#052e16"/><circle cx="152" cy="88" r="5" fill="#052e16"/>',
              '<path d="M128 148 v10 M128 158 l-6 8 M128 158 l6 8" stroke="#f87171" stroke-width="3" fill="none"/>']
    return p


SKIN = {
    "kai": "#e8b98a", "mira": "#f0c8a0", "po": "#d9a06b", "lungvolt": "#c99b76",
    "jump": "#efe9dc", "shell": "#7fb069", "sky": "#7dd3fc", "rooster": "#fef3c7",
    "neoncat": "#3b4470", "tusk": "#a9b7c9", "dash": "#8a5a3b", "venom": "#22c55e",
}
SHAPE = {
    "kai": "human", "mira": "human", "po": "human", "lungvolt": "human",
    "jump": "slim", "shell": "round", "sky": "slim", "rooster": "slim",
    "neoncat": "triangle", "tusk": "round", "dash": "long", "venom": "slim",
}


def build_svg(c: dict) -> str:
    a = c["accent"]
    parts = svg_open(a)
    parts += head(SHAPE[c["id"]], SKIN[c["id"]], a)
    parts += features(c["species"] if c["species"] != "human" else "human", c)
    parts += svg_visor(a)
    parts += svg_body(a)
    parts += svg_name(c["en"], f'{c["th"]} · {c["kind"]}', a)
    parts.append("</svg>")
    return "\n".join(parts) + "\n"


def main() -> None:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    roster = []
    for c in ROSTER:
        (OUT_DIR / f'{c["id"]}.svg').write_text(build_svg(c), encoding="utf-8")
        roster.append({
            "id": c["id"], "name_th": c["th"], "name_en": c["en"],
            "species": c["species"], "kind_th": c["kind"], "role": c["role"],
            "faction": c["faction"], "accent": c["accent"],
            "bio_th": c["th_bio"], "bio_en": c["en_bio"],
            "stats": {"spd": c["spd"], "hdl": c["hdl"], "arm": c["arm"],
                      "tec": c["tec"], "chr": c["chr"]},
            "portrait": f"/assets/characters/{c['id']}.svg",
        })
    JSON_OUT.write_text(json.dumps({"game": "PROJECT: NEON DRIVE",
                                    "set": "launch-roster-v1",
                                    "count": len(roster),
                                    "characters": roster},
                                   ensure_ascii=False, indent=2) + "\n",
                        encoding="utf-8")
    print(f"wrote {len(roster)} portraits -> {OUT_DIR}")
    print(f"wrote roster -> {JSON_OUT}")


if __name__ == "__main__":
    main()
