# MoneyBro Backend


Backend API untuk aplikasi manajemen keuangan pribadi **MoneyBro**.

## Tentang

MoneyBro adalah aplikasi untuk membantu mengelola keuangan pribadi dengan fitur:

- 📊 **Expense Tracking** - Catat pengeluaran harian
- 🔄 **Expense Templates** - Template untuk pengeluaran rutin
- 💳 **Installment Management** - Kelola cicilan (kredit, pinjaman)
- 💰 **Debt Tracking** - Catat hutang piutang
- 📧 **Email Reminders** - Notifikasi H-3, H-2, H-1 sebelum jatuh tempo
- 📈 **Dashboard** - Ringkasan keuangan

## Tech Stack

| Technology | Description |
|------------|-------------|
| Go 1.21+ | Backend language |
| GraphQL | API (gqlgen) |
| PostgreSQL | Database |
| Redis | Cache & session |
| JWT | Authentication |
| Resend | Email service |

## API

GraphQL endpoint: `/graphql`

## License

MIT License - Lihat file [LICENSE](LICENSE) untuk detail.
