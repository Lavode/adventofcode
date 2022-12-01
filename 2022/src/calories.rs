use std::path::PathBuf;

use crate::{error::Error, input};

pub fn from_file(path: PathBuf) -> Result<Vec<Vec<u32>>, Error> {
    let content = input::to_string(path)?;
    let lines = content.split("\n");

    // Calories per elf
    let mut elves: Vec<Vec<u32>> = Vec::new();
    elves.push(Vec::new());

    for line in lines {
        // Separator indicating we now get the data of the next elf
        if line == "" {
            elves.push(Vec::new());
        } else {
            let idx = elves.len() - 1;
            let calories: u32 = line
                .parse()
                .map_err(|e| Error::DataError(format!("Invalid calorie input: {}", e)))?;

            elves[idx].push(calories);
        }
    }

    Ok(elves)
}

#[cfg(test)]
mod tests {
    use std::io::Write;

    use tempfile::NamedTempFile;

    use super::*;

    #[test]
    fn from_file_parses_contents() {
        let mut file = NamedTempFile::new().unwrap();
        write!(file, "1\n2\n3\n\n20\n30\n\n50").unwrap();

        let elves = from_file(PathBuf::from(file.path())).unwrap();
        assert_eq!(elves, vec![vec![1, 2, 3], vec![20, 30], vec![50],])
    }

    #[test]
    fn from_file_balks_on_invalid_input() {
        let mut file = NamedTempFile::new().unwrap();
        write!(file, "1\n2\n3\n\nab\n30\n\n50").unwrap();

        from_file(PathBuf::from(file.path())).unwrap_err();
    }
}
