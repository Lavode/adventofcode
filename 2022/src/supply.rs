use std::{collections::VecDeque, str::FromStr};

use crate::error::Error;

/// Collection of stacks of supply crates.
#[derive(Debug)]
pub struct Supplies {
    pub stacks: Vec<VecDeque<char>>,
}

impl FromStr for Supplies {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let mut stacks = Vec::new();

        // We want to fill the stack from the bottom up, and the list of stack IDs is at the very
        // bottom anyway.
        let mut lines = s.lines().rev();

        let stack_count = match lines.next() {
            Some(line) => {
                let ids: Vec<&str> = line.split_whitespace().collect();
                ids.len()
            }
            None => return Err(Error::DataError("Invalid input for supplies".to_string())),
        };
        // Populate list of stacks
        for _ in 0..stack_count {
            stacks.push(VecDeque::<char>::new());
        }

        for line in lines {
            for i in 0..stack_count {
                // The crate of the i-th stack will start at offset i*4, and be three characters
                // long.
                let supply_crate = &line[i * 4..i * 4 + 3];
                if supply_crate == "   " {
                } else {
                    // Format is [X] where X is a single char identifying the crate.
                    let supply_crate = &supply_crate[1..2].chars().next().unwrap();
                    stacks[i].push_front(*supply_crate);
                }
            }
        }

        Ok(Supplies { stacks })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn supplies_from_str() {
        let input = "    [D]    
[N] [C]    
[Z] [M] [P]
 1   2   3";

        let supplies = Supplies::from_str(input).unwrap();

        assert_eq!(supplies.stacks.len(), 3);
        assert_eq!(supplies.stacks[0], vec!['N', 'Z']);
        assert_eq!(supplies.stacks[1], vec!['D', 'C', 'M']);
        assert_eq!(supplies.stacks[2], vec!['P']);
    }
}
