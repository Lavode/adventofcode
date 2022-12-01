use std::{fs, path::PathBuf};

use crate::error::Error;

/// Return the canonical path where the input file for the given day's challenge is stored.
///
/// This does not perform any validations such as whether the file exists, is readable, or contains
/// the correct data.
pub fn input_path(day: u32) -> PathBuf {
    return PathBuf::from(format!("input/day_{:02}.txt", day));
}

/// Read the given input file to a string.
pub fn to_string(path: PathBuf) -> Result<String, Error> {
    fs::read_to_string(path).map_err(|e| Error::IOError(e))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn input_path_returns_correct_file() {
        let path = input_path(2);
        assert_eq!(path, PathBuf::from("input/day_02.txt"));

        let path = input_path(27);
        assert_eq!(path, PathBuf::from("input/day_27.txt"));
    }
}
