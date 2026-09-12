# Client Q&A — Prepared Answers

## “เกมเสร็จไปกี่เปอร์เซ็นต์แล้ว?”

ตอบว่า: **Pre-production foundation พร้อมสำหรับเริ่ม playable vertical slice** ไม่ควรตอบเป็นเปอร์เซ็นต์รวม เพราะ production game มี art, client, backend, content, QA, certification และ live ops ที่ scale ต่างกันมาก

สิ่งที่เสร็จแล้วคือ world/story structure, 100-quest graph, 25 characters, factions/districts, vertical-slice acceptance criteria, schemas, server-authority rules, tests และ production architecture direction.

## “ทำไมต้อง Unreal?”

เพราะ product ต้องการ 3D vehicle experience, cinematic/environment tooling และ dedicated authoritative multiplayer path ใน engine เดียวกัน ขณะที่ durable progression/economy แยกไป service plane เพื่อความปลอดภัยและการ recover.

## “ทำ MMO ได้จริงไหม?”

สถาปัตยกรรมถูกออกแบบให้ขยายไป persistent multiplayer ได้ แต่ concurrency จริงต้องพิสูจน์ด้วย alpha/load test. ไม่ควรรับประกันตัวเลขผู้เล่นพร้อมกันก่อน benchmark.

## “VIP จะ Pay-to-Win ไหม?”

ไม่. Contract ปัจจุบันอนุญาต capacity/convenience เช่น garage slots, storage และ saved configurations แต่ competitive build signature และ ranked progression ไม่ขึ้นกับ VIP.

## “ทำไม 1 User 1 รถ?”

เพื่อสร้าง emotional attachment และ identity. ผู้เล่นสามารถ rebuild รถคันเดิมได้หลาย role. VIP เพิ่ม collection capacity ภายหลังโดยไม่ทำให้รถแรกหมดความหมาย.

## “ทำไมต้องทำ Vertical Slice ก่อน?”

เพราะมัน lock สิ่งที่แพงที่สุดก่อน: handling feel, visual bar, environment density, quest pacing, networking boundary, persistence และ production speed. ถ้าสิ่งเหล่านี้ไม่ผ่าน การทำทั้งเมืองก่อนจะเพิ่มความเสี่ยง.

## “พรุ่งนี้มีอะไรให้ดู?”

มี interactive browser demo ที่เปิด offline ได้, repository evidence, Game Design Bible v0.2, vertical-slice specification, executable reference tests, CI และ architecture ADRs.

## “Timeline เท่าไร?”

สำหรับ playable vertical slice ใช้ planning range 8–12 สัปดาห์หลัง lock scope. Online alpha 4–6 เดือนโดยประมาณ. Full production เป็น multi-phase program และต้อง estimate หลังตกลง platform/fidelity/content/concurrency.

## “งบเท่าไร?”

ไม่ควรล็อกตัวเลขก่อน scope. ให้เสนอ discovery/vertical-slice package ก่อน แล้วคำนวณ quotation จาก team-months, asset fidelity, platform target, outsourcing, hosting, voice/cinematic และ ownership requirements.
