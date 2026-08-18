# Backend Challenge (Golang)

Repo นี้เป็นโครงสำหรับทำโจทย์จาก [7-solutions/backend-challenge](https://github.com/7-solutions/backend-challenge)

## What to complete

### 1) User Management API (implement)
- Golang + MongoDB + JWT (HS256)
- CRUD users
- Logging middleware (method, path, execution time)
- Background goroutine ทุก 10 วินาที log จำนวน users
- Unit tests (mock MongoDB ตามความเหมาะสม)

### 2) Lottery Search System (design only)
- ออกแบบระบบค้นหาตั๋วล็อตเตอรี่ด้วย wildcard `*`
- ห้าม assign ซ้ำ “ในเวลาเดียวกัน” สำหรับ pattern เดิม
- อธิบาย storage/DB, indexing, algorithm, performance, concurrency strategy

## Notes

- ยังไม่ใส่โค้ดในรอบนี้ — สร้างโครงโปรเจกต์ไว้ก่อน

