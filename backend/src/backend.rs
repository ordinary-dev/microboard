//! Operations with the database and internal parameters.

pub use boards::Board;
pub use migrations::apply_migrations;
pub use posts::Post;
pub use threads::Thread;
pub use users::User;

mod boards;
mod jwt;
mod migrations;
mod posts;
mod threads;
mod users;
