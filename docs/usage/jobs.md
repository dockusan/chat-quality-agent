# Tạo công việc

Công việc (Job) là đơn vị chính để CQA phân tích chat. Có **2 loại công việc**, mỗi loại phục vụ mục đích khác nhau.

## 2 loại công việc

### 1. Đánh giá chất lượng CSKH (QC Analysis)

**Mục đích**: Kiểm tra nhân viên có tuân thủ quy định CSKH không.

**Cách hoạt động**: AI đọc cuộc chat, đối chiếu với bộ quy tắc bạn đặt ra, rồi cho kết quả:
- **Đạt / Không đạt** cho từng cuộc chat
- **Điểm số 0-100**
- **Danh sách vấn đề** (vi phạm cụ thể, bằng chứng, mức độ nghiêm trọng)

**Ví dụ quy tắc QC**:
- Nhân viên phải chào hỏi lịch sự
- Trả lời trong vòng 5 phút
- Không được dùng từ ngữ thiếu chuyên nghiệp
- Phải hỏi khách còn cần hỗ trợ gì không

### 2. Phân loại chat (Classification)

**Mục đích**: Tự động gán nhãn cho cuộc chat theo chủ đề.

**Cách hoạt động**: AI đọc cuộc chat, xác định nội dung thuộc nhãn nào, rồi gán nhãn kèm mô tả ngắn.

**Ví dụ nhãn phân loại**:
- Khiếu nại (khách phàn nàn chất lượng, thái độ...)
- Góp ý (khách đề xuất cải thiện)
- Hỏi giá / Hỏi menu
- Đặt bàn / Đặt hàng
- Hỏi khuyến mãi

---

## Tạo công việc mới

Vào menu **Công việc** ở sidebar, bấm **Tạo công việc**. Wizard gồm 6 bước:

![Tạo công việc](/screenshots/tao-cong-viec.png)

### Bước 1: Thông tin cơ bản

- **Tên công việc** (bắt buộc, tối thiểu 2 ký tự): Ví dụ "QC CSKH hàng ngày", "Phân loại chat"
- **Mô tả** (không bắt buộc): Ghi chú về mục đích
- **Chọn loại công việc**: Đánh giá chất lượng hoặc Phân loại

### Bước 2: Chọn kênh đầu vào

Chọn 1 hoặc nhiều kênh chat đã kết nối. Chỉ cuộc chat từ các kênh được chọn mới được phân tích.

### Bước 3: Cấu hình quy tắc

**Với QC Analysis:**

Viết bộ quy tắc CSKH bằng markdown. Mỗi quy tắc gồm:
- Tên quy tắc
- Mô tả chi tiết
- Mức độ vi phạm

Bấm **Dùng mẫu quy định CSKH** để tải mẫu có sẵn, sau đó chỉnh sửa theo quy định riêng của bạn.

**Điều kiện bỏ qua (Skip)**: Định nghĩa khi nào cuộc chat nên SKIP (không đánh giá):
- Không có tin nhắn trả lời từ nhân viên
- Cuộc chat dưới 2 tin nhắn
- Chỉ có sticker/file, không có nội dung text
- Tin nhắn spam

Bấm **Dùng mẫu điều kiện** để tải mẫu.

**Với Classification:**

Thêm các nhãn phân loại. Mỗi nhãn gồm:
- **Tên nhãn** (bắt buộc): Ví dụ "Khiếu nại"
- **Mô tả** (bắt buộc): Mô tả khi nào chat thuộc nhãn này. Ví dụ "Khách phàn nàn về chất lượng sản phẩm, thái độ phục vụ, hoặc thời gian chờ"
- **Mức độ**: Nghiêm trọng hoặc Cần cải thiện

Bấm **Thêm nhãn** để thêm nhiều nhãn.

### Bước 4: Thông báo đầu ra

Cấu hình nơi nhận kết quả phân tích. Có thể thêm nhiều kênh thông báo:

**Telegram:**
- **Bot Token**: Tạo bot tại [@BotFather](https://t.me/BotFather), lấy token
- **Group ID**: Thêm bot vào group, gửi tin nhắn, vào [@RawDataBot](https://t.me/RawDataBot) để lấy Group ID (số âm)

**Email:**
- **SMTP Host**: Ví dụ `smtp.gmail.com`
- **SMTP Port**: `587` (TLS) hoặc `465` (SSL)
- **Username / Password**: Tài khoản SMTP (với Gmail dùng App Password)
- **Email gửi**: Địa chỉ email gửi đi
- **Email nhận**: Có thể nhập nhiều email, cách nhau bằng dấu phẩy

**Tùy chỉnh nội dung thông báo:**
- Mặc định: Tự động tạo báo cáo tổng hợp + danh sách vấn đề
- Tùy chỉnh: Viết template riêng với các biến: `{{job_name}}`, `{{total}}`, `{{passed}}`, `{{failed}}`, `{{issues}}`, `{{content}}`, `{{link}}`

Bấm **Gửi thử** để kiểm tra — phải thấy tin nhắn test trước khi tiếp tục.

### Bước 5: Lịch chạy

**Lịch phân tích** (khi nào chạy):
- Dùng cron expression (5 trường): phút giờ ngày-tháng tháng ngày-tuần
- Mặc định: `0 7 * * *` (7 giờ sáng mỗi ngày)
- Ví dụ khác:
  - `0 9 * * 1-5` — 9h sáng thứ 2 đến thứ 6
  - `0 */6 * * *` — Mỗi 6 giờ
  - `0 18 * * *` — 6 giờ chiều mỗi ngày

**Lịch gửi kết quả**:
- **Ngay lập tức**: Gửi ngay khi phân tích xong
- **Theo lịch**: Gửi vào giờ nhất định (cron expression)
- **Một lần**: Gửi 1 lần vào thời điểm cụ thể

### Bước 6: Xác nhận

Xem lại toàn bộ cấu hình trước khi tạo. Bấm **Tạo công việc**.

---

## Chạy công việc

### Chạy thủ công

Vào chi tiết công việc, bấm **Chạy ngay**. Có 3 chế độ:

| Chế độ | Mô tả |
|--------|-------|
| **Chưa phân tích** | Chạy tất cả cuộc chat chưa được phân tích bởi công việc này |
| **Từ lần cuối** | Chạy từ cuộc chat cuối cùng đã phân tích trở đi |
| **Tùy chọn** | Chọn khoảng thời gian và/hoặc giới hạn số lượng |

### Chạy thử (Test Run)

Bấm **Chạy thử** — hệ thống chạy trên 3 cuộc chat mẫu để kiểm tra quy tắc và cấu hình AI trước khi chạy thật.

### Chạy tự động

Nếu đã cấu hình lịch chạy (cron), công việc sẽ tự động chạy theo lịch.

---

## Chỉnh sửa công việc

Bấm icon bút chì trên trang chi tiết công việc để chỉnh sửa:
- Tên, mô tả
- Kênh đầu vào
- Quy tắc đánh giá / nhãn phân loại
- Kênh thông báo
- Lịch chạy
