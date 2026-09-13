# Client Presentation Talk Track — Thai

## 1. เปิดการนำเสนอ

“โปรเจกต์นี้ไม่ได้เริ่มจากคำถามว่าเราจะทำเกมแข่งรถอีกเกมหนึ่งอย่างไร แต่เริ่มจากคำถามว่า ถ้ารถในเกมมีคุณค่าทางอารมณ์เท่ากับตัวละคร RPG จะเกิดอะไรขึ้น”

ประโยคหลักของเกมคือ:

> **You don't own the road. You earn it.**

## 2. อธิบายภาพใหญ่

PROJECT: NEON DRIVE เป็น Online Open World RPG ใน NOVA CITY ปี 2097

สิ่งที่ผู้เล่นทำไม่ได้มีแค่แข่งรถ แต่ประกอบด้วย:
- สร้างรถ,
- ทำงาน,
- หาอะไหล่,
- สร้างชื่อเสียง,
- รู้จัก NPC,
- เลือก faction,
- ตั้ง crew,
- แข่งขัน,
- เปิดเผยความลับของ DRIVE ZERO.

## 3. จุดขายที่ต้องย้ำ

“รถไม่ใช่ Item ที่ซื้อแล้วทิ้ง”

รถหนึ่งคันมี:
- Vehicle ID,
- เจ้าของ,
- build revision,
- ประวัติการซ่อม,
- ประวัติการแข่งขัน,
- reputation,
- provenance.

เพราะฉะนั้นรถคันแรกที่ดูเหมือนเศษเหล็กสามารถกลายเป็นรถ Legendary ที่ทั้ง server รู้จักได้

## 4. Story Hook

พาผู้เล่นไปที่ Garage 17 และ Maya Voss

Maya มอบรถเก่าให้ผู้เล่น แต่ต่อมาพบว่ารถนี้เป็น prototype ของ PROJECT DRIVE ZERO และ ECU มีสิ่งที่ VANTEX ตามหามาหลายปี

ตรงนี้ทำให้ “ระบบสร้างรถ” กับ “เนื้อเรื่องหลัก” เป็นเรื่องเดียวกัน ไม่ใช่ระบบที่แยกกันอยู่

## 5. แสดง Interactive Demo

ลำดับแนะนำ:

1. Hero / metrics — 100 quests, 10 districts, 25 characters, 5 factions
2. NOVA CITY — กดดู district
3. Story — 7 chapters
4. Vehicle Lab — เปลี่ยน Street / Drift / Circuit / Off-road
5. Vertical Slice — กด Next จาก MQ001 ไป First Ignition
6. VIP — ชี้ให้เห็นว่าเพิ่ม garage ไม่เพิ่ม power
7. Architecture — อธิบาย server authority
8. Roadmap — ปิดด้วยสิ่งที่จะส่งมอบต่อ

## 6. Technology

สำหรับ production direction เราเลือก:
- Unreal Engine 5.8 สำหรับ client และ dedicated gameplay server,
- Go service plane,
- PostgreSQL,
- Redis,
- NATS เฉพาะ workflow ที่ต้องใช้ event durability,
- OpenTelemetry observability,
- containerization ก่อน และ Kubernetes เมื่อ load evidence บอกว่าจำเป็น.

ประโยคสำคัญ:

“เราไม่เอา microservices หรือ Kubernetes มาเป็นจุดขาย เราเอามาใช้เมื่อมันแก้ปัญหาที่พิสูจน์แล้วว่ามีจริง”

## 7. สิ่งที่มีอยู่แล้ว

ห้ามพูดว่า “เกมเสร็จแล้ว”

ให้พูดว่า:

“Pre-production foundation และ system contracts เสร็จแล้ว และมี executable reference tests สำหรับ ownership, idempotency, vehicle build revision และ race result binding พร้อม CI”

## 8. เสนอขั้นต่อไป

แนะนำขาย **Playable Vertical Slice** ก่อน

เหตุผล:
- ลูกค้าเห็น product feel จริง,
- lock visual quality,
- lock handling,
- test story pacing,
- test networking boundaries,
- ใช้เป็น milestone สำหรับ investor/publisher/internal approval ได้.

## 9. คำถามปิดการขาย

ถามลูกค้าตรง ๆ:

1. อยากเห็นเกมบน PC ก่อน หรือมี Console ใน milestone แรก?
2. ต้องการ realistic AAA look หรือ stylized/AA ที่ production เร็วกว่า?
3. Vertical Slice ต้องรองรับกี่ผู้เล่นพร้อมกัน?
4. ลูกค้าต้องการเป็นเจ้าของ IP/source ทั้งหมด หรือทำในรูปแบบ co-development?
5. เป้าหมายถัดไปคือ investor demo, publisher demo, internal greenlight หรือ commercial alpha?

## 10. ประโยคปิด

“สิ่งที่เราขายวันนี้ไม่ใช่แค่ idea แต่คือ world, content structure, technical authority model และ production path ที่สามารถเริ่มสร้าง vertical slice ได้ทันที”
