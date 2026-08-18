# Postman — Backend Challenge

Collection และ environment สำหรับยิง API ที่มีตอนนี้: health, register, login และทดสอบ JWT บน `/users`

## Import

1. เปิด Postman → **Import**
2. Import ทั้งสองไฟล์:
   - [`Backend_Challenge.postman_collection.json`](./Backend_Challenge.postman_collection.json)
   - [`Backend_Challenge_Local.postman_environment.json`](./Backend_Challenge_Local.postman_environment.json)
3. เลือก environment **Backend Challenge Local** มุมขวาบน
4. `base_url` ต้องเป็น `http://localhost:8080` ให้ตรงกับ `PORT` ใน [`.env.example`](../../.env.example)

## Prerequisites

| โหมด | คำสั่ง | หมายเหตุ |
|------|--------|----------|
| Mongo | `docker compose up -d` | ต้องเปิด Docker Desktop ก่อน |
| API | จากรากโปรเจกต์ `go run ./cmd/api` | อ่าน `.env` แล้วต่อ Mongo จริง |

## Auth model

| Route | JWT |
|-------|-----|
| `GET /health` | ไม่ใช้ |
| `POST /auth/register` | ไม่ใช้ |
| `POST /auth/login` | ไม่ใช้ (ได้ token กลับมา) |
| `/users/*` | ใช้ `Authorization: Bearer {{access_token}}` |

CRUD `/users` ยังไม่มี handler รอบนี้ collection จึงมีเคส JWT บน `GET /users` ไว้ก่อน

หมายเหตุ: ถ้ายังไม่มี route ในกลุ่ม `/users` Fiber อาจตอบ **404** แทน 401 เพราะ middleware ของ group ทำงานเมื่อมี route ตรงกับ path นั้น หลังมี CRUD เคส 401 จะตรงตามที่เขียนไว้

## Happy path order

1. **Health → GET Health - OK**
2. **Auth → Register - Happy** — Pre-request จะสุ่ม email ให้ไม่ชน unique index
3. **Auth → Login - Happy** — Tests เก็บ `access_token`

## Negative cases included

| Request | Expect |
|---------|--------|
| Register - Invalid JSON | 400 `code: 400` |
| Register - Missing Fields | 400 `code: 400` |
| Register - Email Conflict | 409 `code: 409` (ยิงหลัง Happy) |
| Login - Missing Fields | 400 `code: 400` |
| Login - Invalid Credentials | 401 `code: 401` |
| Users - Missing JWT | 401 `code: 401` |
| Users - Invalid JWT | 401 `code: 401` |

## Environment keys

| Key | Purpose |
|-----|---------|
| `base_url` | API host |
| `register_name` / `register_email` / `register_password` | body ของ register/login |
| `access_token` | ถูกเติมหลัง login สำเร็จ |
| `user_id` | ถูกเติมหลัง register สำเร็จ |

## Envelope

| code | Meaning |
|------|---------|
| `0` | success |
| `400` | invalid parameter |
| `401` | unauthorized |
| `404` | not found |
| `409` | conflict |
| `500` | internal error |

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "serverTime": "2026-08-18T19:01:31+07:00"
}
```

`GET /health` คืน `{ "status": "ok" }` ไม่ใช้ envelope
