# Kết nối MCP

MCP (Model Context Protocol) cho phép Claude Web hoặc Claude Desktop truy vấn dữ liệu CQA trực tiếp. Bạn có thể hỏi Claude về cuộc chat, kết quả đánh giá, thống kê... mà không cần mở CQA.

## MCP là gì?

MCP là giao thức kết nối AI với các hệ thống bên ngoài. Khi kết nối CQA với Claude qua MCP, bạn có thể:

- "Hôm nay có bao nhiêu cuộc chat khiếu nại?"
- "Tóm tắt vấn đề CSKH tuần này"
- "Cuộc chat nào bị điểm thấp nhất hôm nay?"
- "Nhân viên nào bị nhiều vi phạm nhất?"

Claude sẽ tự truy vấn CQA và trả lời.

## Tạo kết nối MCP

Vào menu **MCP** ở sidebar, bấm **Tạo kết nối**.

1. Nhập **Tên kết nối** (ví dụ "Claude Desktop", "Claude Web")
2. Bấm **Tạo**
3. Hệ thống trả về:
   - **Client ID**: Mã định danh kết nối
   - **Client Secret**: Khóa bí mật (**chỉ hiển thị 1 lần**, copy ngay!)

::: danger Quan trọng
**Client Secret chỉ hiển thị 1 lần** khi tạo. Nếu quên copy, bạn phải xóa kết nối và tạo lại.
:::

## Kết nối với Claude Desktop

Mở file cấu hình Claude Desktop (`claude_desktop_config.json`):

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Thêm cấu hình MCP server:

```json
{
  "mcpServers": {
    "cqa": {
      "url": "https://cqa.yourdomain.com/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_CLIENT_SECRET"
      }
    }
  }
}
```

Thay `cqa.yourdomain.com` bằng URL CQA và `YOUR_CLIENT_SECRET` bằng secret vừa copy.

Khởi động lại Claude Desktop. Bạn sẽ thấy icon CQA trong danh sách MCP tools.

## Kết nối với Claude Web

Claude Web hỗ trợ MCP qua OAuth:

1. Vào Claude Web ([claude.ai](https://claude.ai))
2. Bấm icon MCP > **Add MCP Server**
3. Nhập URL: `https://cqa.yourdomain.com/mcp`
4. Xác thực bằng Client ID và Client Secret

## Các công cụ MCP có sẵn

| Tool | Mô tả |
|------|-------|
| **conversations** | Lấy danh sách cuộc hội thoại gần đây |
| **transcripts** | Đọc nội dung tin nhắn của 1 cuộc chat |
| **evaluations** | Xem kết quả đánh giá QC |
| **statistics** | Thống kê tổng quan (số chat, tỉ lệ đạt, chi phí...) |

## Thu hồi kết nối

Bấm **Thu hồi** bên cạnh kết nối trong danh sách. Sau khi thu hồi, Claude không thể truy cập CQA qua kết nối đó nữa.
