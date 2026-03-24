# FAQ & Xử lý lỗi

## Câu hỏi thường gặp

### CQA hỗ trợ những kênh chat nào?

Hiện tại hỗ trợ **Zalo OA** và **Facebook Messenger**. Các kênh khác (Viber, Shopee Chat...) sẽ được bổ sung trong tương lai.

### Nên dùng Claude hay Gemini?

| Tiêu chí | Claude | Gemini |
|----------|--------|--------|
| Chất lượng tiếng Việt | Tốt hơn | Khá |
| Chi phí | Cao hơn | Rẻ hơn nhiều |
| Tốc độ | Nhanh | Nhanh hơn |
| Khuyến nghị QC | Sonnet 4.6 | Pro 2.5 |
| Khuyến nghị Phân loại | Haiku 4.5 | Flash 2.0 |

**Kết luận**: Claude cho đánh giá QC chính xác hơn. Gemini cho phân loại đơn giản và tiết kiệm.

### Batch mode là gì? Có nên bật không?

Batch mode gom nhiều cuộc chat vào 1 lần gọi AI, tiết kiệm 60-80% chi phí token. **Nên bật** với batch size 5-10. Xem chi tiết tại [Cấu hình AI](/usage/ai-settings).

### Tại sao Zalo OA chỉ lấy được tin nhắn 48 giờ?

Do giới hạn API của Zalo. Zalo chỉ cho phép ứng dụng bên thứ 3 đọc tin nhắn trong cửa sổ 48 giờ gần nhất. Tin nhắn cũ hơn không truy cập được.

**Giải pháp**: Đặt lịch đồng bộ thường xuyên (15-30 phút/lần) để không bỏ sót tin nhắn.

### CQA cần bao nhiêu tài nguyên?

- **App**: ~100MB RAM
- **MySQL**: ~500MB RAM
- **Nginx**: ~50MB RAM
- **Disk**: Tùy số lượng tin nhắn, thường dưới 5GB cho 100K cuộc chat
- **Tổng**: VPS 1GB RAM + 10GB disk là đủ dùng

### Có cần tên miền không?

Không bắt buộc. CQA chạy được với IP trực tiếp (http://IP-VPS). Tên miền + SSL chỉ cần khi muốn HTTPS hoặc dùng MCP với Claude Web.

### Chi phí AI ước tính bao nhiêu?

Ví dụ đánh giá 100 cuộc chat/ngày:
- Claude Sonnet + Batch 5: ~$0.80/ngày (~600K VND/tháng)
- Claude Haiku + Batch 10: ~$0.15/ngày (~120K VND/tháng)
- Gemini Flash + Batch 10: ~$0.03/ngày (~24K VND/tháng)

### Dữ liệu có an toàn không?

- Mật khẩu được hash (bcrypt)
- API key và credential kênh chat được mã hóa AES-256
- JWT token với refresh token rotation
- Khóa tài khoản sau 5 lần đăng nhập sai (15 phút)
- HTTPS (nếu bật SSL)

---

## Xử lý lỗi

### Không truy cập được CQA sau cài đặt

```bash
# Kiểm tra container
docker compose ps

# Xem log
docker compose logs --tail=20
```

**Nguyên nhân phổ biến:**
- Port 80 bị chặn bởi firewall → Mở port: `ufw allow 80`
- Container chưa start → Chờ 30 giây rồi kiểm tra lại
- Image sai kiến trúc → Xem log có `exec format error` không

### exec format error

Image Docker sai kiến trúc (ví dụ ARM image trên AMD64 server).

```bash
docker compose down
docker rmi buitanviet/chat-quality-agent:latest buitanviet/chat-quality-agent-nginx:latest
docker compose pull
docker compose up -d
```

### Kênh Zalo báo lỗi xác thực

Token Zalo hết hạn (90 ngày). Bấm **Xác thực lại** trên trang kênh chat.

### Kênh Facebook không đồng bộ được

- Kiểm tra Page Access Token còn hiệu lực
- Kiểm tra Page ID đúng
- Đảm bảo token có quyền `pages_messaging`

### AI trả kết quả không chính xác

1. **Kiểm tra quy tắc** — Quy tắc quá chung chung sẽ cho kết quả không rõ ràng. Viết càng chi tiết càng tốt.
2. **Giảm batch size** — Batch size lớn giảm độ chính xác. Thử giảm về 3-5.
3. **Đổi model** — Dùng model mạnh hơn (Claude Sonnet thay Haiku, Gemini Pro thay Flash).
4. **Chạy thử** — Dùng nút "Chạy thử" để test trên 3 cuộc chat trước khi chạy thật.

### Thông báo Telegram không gửi được

- Bot Token đúng chưa? Test bằng: `https://api.telegram.org/bot<TOKEN>/getMe`
- Bot đã được thêm vào group chưa?
- Group ID đúng chưa? (phải là số âm)
- Bot có quyền gửi tin nhắn trong group không?

### Thông báo Email không gửi được

- SMTP host và port đúng chưa?
- Username/password SMTP đúng chưa?
- Với Gmail: Đã dùng App Password chưa? (không dùng mật khẩu Gmail thường)
- Port 587 (TLS) hay 465 (SSL) — thử cả 2

### SSL không hoạt động

```bash
docker compose logs nginx --tail=20
```

**Nguyên nhân phổ biến:**
- DNS chưa trỏ đúng IP → `ping cqa.yourdomain.com` kiểm tra
- Port 80 bị chặn (Let's Encrypt cần port 80 để xác minh) → Mở port 80
- Đã request quá 5 lần/tuần → Chờ 1 tuần

### Quên mật khẩu admin

Nếu là admin duy nhất và quên mật khẩu, cần reset trực tiếp trong database:

```bash
cd /opt/cqa
docker compose exec db mysql -u root -p$MYSQL_ROOT_PASSWORD cqa

# Trong MySQL:
UPDATE users SET password_hash = '$2a$10$...' WHERE email = 'admin@example.com';
```

Tốt hơn: Thêm admin mới qua API hoặc liên hệ người có quyền Owner để reset.
