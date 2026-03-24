# Người dùng & phân quyền

CQA hỗ trợ phân quyền theo mô hình **Owner > Admin > Member**.

## Vai trò (Role)

| Vai trò | Quyền hạn |
|---------|-----------|
| **Owner** | Toàn quyền. Tạo/xóa công ty, quản lý tất cả user (kể cả admin). Mỗi công ty có 1 owner. |
| **Admin** | Quản lý kênh, công việc, xem kết quả, mời/xóa member. Không thể xóa công ty hoặc quản lý admin khác. |
| **Member** | Quyền hạn chế theo cấu hình. Chỉ truy cập các tính năng được cho phép. |

## Thêm người dùng

Vào menu **Người dùng** ở sidebar, bấm **Thêm người dùng**.

![Quản lý người dùng](/screenshots/quan-ly-user.png)

### Thông tin cần nhập

| Trường | Bắt buộc | Mô tả |
|--------|----------|-------|
| **Tên hiển thị** | Có | Tên hiển thị trên giao diện |
| **Email** | Có | Email đăng nhập (phải là email hợp lệ) |
| **Mật khẩu** | Có | Tối thiểu 8 ký tự, có chữ hoa và số |
| **Vai trò** | Có | Admin hoặc Member |

## Phân quyền cho Member

Khi chọn vai trò **Member**, bạn sẽ thấy bảng phân quyền:

| Tính năng | View (Xem) | Edit (Sửa) |
|-----------|-----------|-----------|
| **Kênh chat** | Xem danh sách kênh | Thêm/sửa/xóa kênh, đồng bộ |
| **Tin nhắn** | Xem cuộc hội thoại | Xuất tin nhắn |
| **Công việc** | Xem danh sách, kết quả | Tạo/sửa/chạy công việc |
| **Cài đặt** | Xem cấu hình | Thay đổi cấu hình |

- **Không tick**: Không truy cập được tính năng
- **Tick View**: Chỉ xem, không sửa
- **Tick Edit**: Xem và sửa (bao gồm quyền View)

## Đổi vai trò

Trong danh sách người dùng, thay đổi vai trò qua dropdown **Role**:
- Chuyển Member → Admin: user có toàn quyền (không cần cấu hình phân quyền)
- Chuyển Admin → Member: cần cấu hình phân quyền chi tiết

::: warning Lưu ý
Admin chỉ có thể quản lý Member. Để thay đổi vai trò Admin khác, cần quyền Owner.
:::

## Đặt lại mật khẩu

Admin/Owner có thể reset mật khẩu cho user khác:

1. Trong danh sách người dùng, bấm icon **khóa** bên cạnh user
2. Nhập mật khẩu mới (tối thiểu 8 ký tự, có chữ hoa và số)
3. Bấm **Đặt lại**

User sẽ bị đăng xuất khỏi tất cả thiết bị và phải đăng nhập lại bằng mật khẩu mới.

## Đổi mật khẩu của chính mình

Vào menu **Cài đặt** > tab **Đổi mật khẩu**:
1. Nhập mật khẩu hiện tại
2. Nhập mật khẩu mới
3. Bấm **Đổi mật khẩu**

## Xóa người dùng

Bấm icon thùng rác bên cạnh user. Xóa user khỏi công ty (user vẫn tồn tại trong hệ thống, có thể được mời lại).
