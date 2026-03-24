# Xem kết quả

Sau khi công việc chạy xong, vào menu **Công việc** > bấm vào công việc cần xem.

## Thống kê tổng quan (KPI Cards)

### Kết quả QC Analysis

| Chỉ số | Ý nghĩa |
|--------|---------|
| **Tổng cuộc chat** | Số cuộc chat đã phân tích |
| **Tỉ lệ đạt** | Phần trăm cuộc chat đạt yêu cầu |
| **Vấn đề phát hiện** | Tổng số vi phạm |
| **Điểm trung bình** | Điểm trung bình 0-100 |

### Kết quả Phân loại

| Chỉ số | Ý nghĩa |
|--------|---------|
| **Tổng cuộc chat** | Số cuộc chat đã xử lý |
| **Đã phân loại** | Số cuộc chat được gán nhãn (trừ SKIP) |
| **Bỏ qua** | Số cuộc chat bị SKIP (không đủ nội dung) |

## Chế độ xem kết quả

### Tab Kết quả đánh giá

![Kết quả đánh giá QC](/screenshots/ket-qua-cong-viec-danh-gia.png)

![Kết quả phân loại](/screenshots/ket-qua-cong-viec-phan-loai.png)

Danh sách kết quả có 2 chế độ hiển thị:

**Chế độ danh sách** (mặc định): Mỗi cuộc chat hiển thị trên 1 dòng:
- Tên khách hàng
- Thời gian chat
- Trạng thái (Đạt/Không đạt/Bỏ qua)
- Nhãn phân loại (với công việc phân loại)
- Số vấn đề phát hiện
- Bấm mũi tên mở rộng để xem chi tiết

**Chế độ bảng**: Hiển thị dạng bảng, dễ so sánh nhiều kết quả.

### Xem chi tiết 1 kết quả

Bấm mở rộng 1 dòng, bạn sẽ thấy:

**Với QC Analysis:**
- **Diễn biến cuộc chat**: Toàn bộ tin nhắn giữa khách và nhân viên
- **Đánh giá chi tiết**: Kết quả Đạt/Không đạt, điểm số
- **Danh sách vấn đề**: Từng vi phạm với mức độ (Nghiêm trọng/Cần cải thiện), tên quy tắc, bằng chứng cụ thể

**Với Phân loại:**
- **Diễn biến cuộc chat**: Nội dung hội thoại
- **Kết quả phân loại**: Nhãn được gán (ví dụ "Khiếu nại"), mô tả ngắn (ví dụ "Khách phàn nàn về chất lượng đồ uống và thời gian phục vụ chậm")

### Bộ lọc kết quả

| Bộ lọc | Áp dụng cho | Mô tả |
|--------|-------------|-------|
| **Đã phân loại / Tất cả / Bỏ qua** | Phân loại | Lọc theo trạng thái |
| **Lọc loại** | Phân loại | Lọc theo nhãn (Khiếu nại, Góp ý...) |
| **Đạt / Không đạt** | QC | Lọc theo kết quả |

### Tab Lịch sử chạy

Danh sách các lần chạy công việc:
- Thời gian chạy
- Trạng thái (Thành công / Lỗi / Đang chạy)
- Số cuộc chat: tổng / đã phân tích / đạt / bỏ qua
- Thời gian chạy
- Chi phí (số token, USD)
- Lỗi (nếu có)

## Xuất kết quả

Bấm nút **CSV** hoặc **Excel** phía trên danh sách kết quả.

### Nội dung file xuất

**QC Analysis (CSV/Excel):**

| Cột | Ví dụ |
|-----|-------|
| Tên khách | Nguyễn Văn A |
| Ngày phát sinh chat | 22/03/2026 |
| Ngày đánh giá | 23/03/2026 |
| Kết quả | Không đạt |
| Đánh giá | Nhân viên không chào hỏi, trả lời chậm |
| Điểm | 45 |
| Vấn đề | Thiếu lời chào; Phản hồi chậm 20 phút |

**Phân loại (CSV/Excel):**

| Cột | Ví dụ |
|-----|-------|
| Tên khách | Trần Thị B |
| Ngày phát sinh chat | 22/03/2026 |
| Ngày đánh giá | 23/03/2026 |
| Loại | Khiếu nại |
| Vấn đề | Khách phàn nàn về chất lượng đồ uống |
| Nội dung chat | (tóm tắt nội dung) |

## Xóa kết quả

Bấm **Xóa kết quả** (icon thùng rác đỏ) để xóa toàn bộ kết quả phân tích của công việc. Thao tác này không thể hoàn tác.
