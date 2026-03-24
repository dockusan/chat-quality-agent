# Giới thiệu

**Chat Quality Agent (CQA)** la he thong ma nguon mo giup doanh nghiep tu dong phan tich chat luong cham soc khach hang (CSKH) qua cac kenh chat.

## Vấn đề CQA giải quyết

Khi doanh nghiệp có nhiều kênh chat (Zalo OA, Facebook Messenger), việc kiểm tra chất lượng CSKH thủ công tốn nhiều thời gian và dễ bỏ sót. CQA tự động hóa quy trình này:

- **Đọc hết mọi cuộc chat** — Đồng bộ tự động từ Zalo OA và Facebook Messenger
- **Đánh giá bằng AI** — AI đọc cuộc chat, chấm điểm 0-100, phát hiện vi phạm theo quy định CSKH của bạn
- **Phân loại tự động** — Gán nhãn cho cuộc chat: khiếu nại, góp ý, hỏi giá, đặt bàn...
- **Cảnh báo ngay** — Gửi thông báo qua Telegram hoặc Email khi phát hiện vấn đề

## Tính năng chính

| Tính năng | Mô tả |
|-----------|-------|
| Đồng bộ tin nhắn | Tự động lấy tin nhắn từ Zalo OA và Facebook Messenger |
| Đánh giá CSKH (QC) | AI chấm điểm, phân loại Đạt/Không đạt, chỉ ra lỗi cụ thể |
| Phân loại chat | Phân loại theo chủ đề tùy chỉnh (khiếu nại, góp ý, hỏi giá...) |
| Cảnh báo tự động | Gửi kết quả qua Telegram và Email theo lịch hẹn |
| Batch AI mode | Gom nhiều cuộc chat/lần gọi AI, tiết kiệm 60-80% chi phí |
| Dashboard | Biểu đồ, thống kê, cảnh báo gần đây |
| Multi-tenant | Nhiều công ty trên 1 hệ thống, phân quyền Owner > Admin > Member |
| Tích hợp MCP | Kết nối với Claude Web/Desktop để truy vấn dữ liệu |
| SSL tự động | Let's Encrypt tự động tạo và gia hạn certificate |
| Xuất dữ liệu | Xuất kết quả ra CSV, Excel; xuất tin nhắn ra TXT, CSV |
| Dữ liệu demo | Import 220 mẫu hội thoại để trải nghiệm trước khi dùng thật |

## Kiến trúc hệ thống

```
+-----------+     +-----------+
| Zalo OA   |---->|           |     +----------+
+-----------+     |  CQA App  |---->| MySQL DB |
+-----------+     |  (Go)     |     +----------+
| Facebook  |---->|           |
+-----------+     +-----+-----+
                        |
                  +-----+-----+
                  |   Nginx   |
                  | (SSL/Proxy)|
                  +-----------+
                        |
              +---------+---------+
              |                   |
        +-----+-----+     +------+------+
        | Claude AI  |     | Gemini AI   |
        +-----+-----+     +------+------+
              |                   |
        +-----+-----+     +------+------+
        | Telegram   |     | Email SMTP  |
        +-----------+     +-------------+
```

## Yêu cầu hệ thống

- **VPS**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)
- **RAM**: Tối thiểu 1GB (khuyến nghị 2GB)
- **Disk**: Tối thiểu 10GB
- **Docker**: Docker Engine 20+ và Docker Compose v2
- **AI API Key**: Claude (Anthropic) hoặc Gemini (Google) — cần ít nhất 1 key

## Bước tiếp theo

- [Cài đặt](/guide/installation) — Cài đặt CQA lên VPS
- [Thiết lập ban đầu](/guide/initial-setup) — Tạo admin và cấu hình
