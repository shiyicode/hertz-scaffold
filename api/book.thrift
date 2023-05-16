// api/book.thrift
namespace go book
namespace py book

struct Book {
    1: string title
    2: string author
    3: i32 pages
}

struct ListBooksRequest {
    3: i32 offset
    4: i32 limit
}
struct ListBooksResponse {
    1: list<Book> books
    2: i32 total
}

struct GetBookRequest {
    1: required i32 id
}

struct CreateBookRequest {
    1: Book book
}

struct UpdateBookRequest {
    1: required Book book
}

struct RenameBookRequest {
    1: required i32 id
    2: required string title
}

struct RenameBookResponse {
    1: required Book book
}

struct DeleteBookRequest {
    1: required i32 id (api.path="id")
}

service BookService {
    ListBooksResponse listBooks(1: ListBooksRequest req) (api.get="/books");
    Book getBook(1: GetBookRequest req) (api.get="/books/:id");
    Book createBook(1: CreateBookRequest req) (api.post="/books/:id");
    Book updateBook(1: UpdateBookRequest req) (api.put="/books/:id");
    Book renameBook(1: RenameBookRequest req) (api.patch="/books/:id");
    Book deleteBook(1: DeleteBookRequest req) (api.delete="/books/:id");
}