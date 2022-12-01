#[derive(Debug)]
pub enum Error {
    /// Error when content of input data is invalid
    DataError(String),
    /// Error when IO operation failed
    IOError(std::io::Error),
}

impl std::error::Error for Error {}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::DataError(s) => write!(f, "Invalid input data: {}", s),
            Self::IOError(e) => write!(f, "IO Error: {e}"),
        }
    }
}
