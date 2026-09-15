#!/usr/bin/env python3
"""Generate two deterministic, low-poly 3D character prototypes as OBJ/MTL."""
from pathlib import Path
import math

OUT = Path(__file__).resolve().parents[1] / "game" / "Content" / "Characters" / "Prototype"

def mesh(name, gender, scale, jacket, accent):
    verts, faces, mats = [], [], []
    def sphere(cx, cy, cz, rx, ry, rz, material, rings=6, seg=10):
        base = len(verts) + 1
        for r in range(rings + 1):
            p = math.pi * r / rings
            for s in range(seg):
                a = 2 * math.pi * s / seg
                verts.append((cx + rx * math.sin(p) * math.cos(a), cy + ry * math.sin(p) * math.sin(a), cz + rz * math.cos(p)))
        for r in range(rings):
            for s in range(seg):
                a, b = base + r * seg + s, base + r * seg + (s + 1) % seg
                c, d = base + (r + 1) * seg + (s + 1) % seg, base + (r + 1) * seg + s
                faces.extend([((a, b, c), material), ((a, c, d), material)])
    def box(cx, cy, cz, rx, ry, rz, material):
        base = len(verts) + 1
        verts.extend([(cx+x, cy+y, cz+z) for x,y,z in [(-rx,-ry,-rz),(rx,-ry,-rz),(rx,ry,-rz),(-rx,ry,-rz),(-rx,-ry,rz),(rx,-ry,rz),(rx,ry,rz),(-rx,ry,rz)]])
        for f in [(1,2,3,4),(5,8,7,6),(1,5,6,2),(2,6,7,3),(3,7,8,4),(5,1,4,8)]: faces.append((tuple(base+i-1 for i in f), material))
    # Z-up human silhouette: boots, legs, torso, head, arms.
    box(0, 0, .12, .18, .28, .12, "boots"); box(0, 0, .52, .18, .24, .28, "pants")
    box(0, 0, 1.12, .34 * scale, .22, .42, jacket); sphere(0, 0, 1.78 * scale, .22 * scale, .20 * scale, .25 * scale, "skin")
    box(-.48 * scale, 0, 1.15 * scale, .12, .16, .38 * scale, jacket); box(.48 * scale, 0, 1.15 * scale, .12, .16, .38 * scale, jacket)
    box(0, -.24, 1.22 * scale, .16 * scale, .035, .08 * scale, accent)
    lines = [f"# {name} — {gender} low-poly prototype", "mtllib " + name + ".mtl", "o " + name]
    lines += [f"v {x:.5f} {y:.5f} {z:.5f}" for x,y,z in verts]
    current = None
    for face, material in faces:
        if material != current: lines.append("usemtl " + material); current = material
        lines.append("f " + " ".join(map(str, face)))
    (OUT / (name + ".obj")).write_text("\n".join(lines) + "\n")
    (OUT / (name + ".mtl")).write_text("\n".join(["newmtl skin\nKd 0.55 0.28 0.18", f"newmtl {jacket}\nKd 0.04 0.12 0.18", "newmtl pants\nKd 0.03 0.04 0.08", "newmtl boots\nKd 0.01 0.01 0.015", f"newmtl {accent}\nKd 0.0 0.9 1.0"]) + "\n")

OUT.mkdir(parents=True, exist_ok=True)
mesh("NovaFemale", "female", .96, "jacket_cyan", "accent_pink")
mesh("RexMale", "male", 1.04, "jacket_magenta", "accent_cyan")
(OUT / "README.md").write_text("""# Prototype 3D Characters\n\nGenerated deterministic OBJ/MTL meshes for Unreal import. NovaFemale and RexMale are low-poly blockout characters for pipeline and gameplay testing; they are not rigged, animated, photorealistic, or final production art. Convert/import through the project's Unreal asset pipeline and replace with skeletal meshes before release.\n""")
