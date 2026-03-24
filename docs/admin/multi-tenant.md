# Quản lý đa công ty

CQA hỗ trợ **multi-tenant** — 1 hệ thống phục vụ nhiều công ty, mỗi công ty có dữ liệu riêng biệt.

## Tạo công ty mới

![Danh sách công ty](/screenshots/danh-sach-cong-ty.png)

Từ trang chính (sau khi đăng nhập), bấm **Thêm công ty**:

1. **Tên công ty** (2-100 ký tự): Ví dụ "SePay Coffee"
2. **Slug** (tự động tạo từ tên): Dùng cho URL, chỉ chứa a-z, 0-9, dấu gạch ngang
   - Ví dụ: "SePay Coffee" → `sepay-coffee`
   - Có thể sửa thủ công nếu muốn
3. Bấm **Lưu**

Người tạo công ty tự động trở thành **Owner** của công ty đó.

## Chuyển đổi giữa các công ty

Bấm icon **mũi tên quay lại** ở góc trên trái sidebar → về trang danh sách công ty → chọn công ty khác.

## Dữ liệu riêng biệt

Mỗi công ty có dữ liệu hoàn toàn tách biệt:
- Kênh chat riêng
- Tin nhắn riêng
- Công việc và kết quả riêng
- Cấu hình AI riêng (provider, key, model)
- Người dùng và phân quyền riêng
- Chi phí AI riêng
- Cài đặt chung riêng (tên, múi giờ, tỉ giá)

## Quản lý người dùng đa công ty

Một user (email) có thể thuộc nhiều công ty với vai trò khác nhau:
- Owner ở công ty A, Member ở công ty B
- Admin ở cả công ty C và D

## Xóa công ty

Chỉ **Owner** mới có quyền xóa công ty.

::: danger Cảnh báo
Xóa công ty sẽ xóa **tất cả dữ liệu** liên quan:
- Tất cả kênh chat
- Tất cả cuộc hội thoại và tin nhắn
- Tất cả công việc và kết quả
- Tất cả chi phí AI, thông báo, activity log
- Liên kết user-công ty (user không bị xóa khỏi hệ thống)

Thao tác này **không thể hoàn tác**.
:::

## Kịch bản sử dụng

**Quản lý chuỗi cửa hàng:**
- Công ty "SePay Coffee HCM" — kết nối OA chi nhánh HCM
- Công ty "SePay Coffee HN" — kết nối OA chi nhánh Hà Nội
- Mỗi chi nhánh có quy tắc CSKH riêng, nhân viên riêng

**Agency quản lý nhiều khách hàng:**
- Công ty "Khách A" — kết nối kênh của khách A
- Công ty "Khách B" — kết nối kênh của khách B
- Mỗi khách có cấu hình AI và quy tắc đánh giá riêng
