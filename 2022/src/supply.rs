use std::{collections::VecDeque, str::FromStr};

use regex::Regex;

use crate::error::Error;

/// Collection of stacks of supply crates.
#[derive(Debug, PartialEq, Eq)]
pub struct Supplies {
    pub stacks: Vec<VecDeque<char>>,
}

impl Supplies {
    pub fn apply(&mut self, command: &Command, command_type: CommandType) -> Result<(), Error> {
        if command.source >= self.stacks.len() {
            return Err(Error::DataError(format!(
                "Invalid command {:?}: Accessing invalid source stack {}",
                command, command.source
            )));
        }

        if command.target >= self.stacks.len() {
            return Err(Error::DataError(format!(
                "Invalid command {:?}: Accessing invalid target stack {}",
                command, command.source
            )));
        }

        // TODO this is a fair bit of repeated code
        match command_type {
            CommandType::SingleCrate => {
                for _ in 0..command.count {
                    let crate_ = match self.stacks[command.source].pop_front() {
                        Some(c) => c,
                        None => {
                            return Err(Error::DataError(format!(
                                "Invalid command {:?}: Popped from empty stack {}",
                                command, command.source
                            )))
                        }
                    };

                    self.stacks[command.target].push_front(crate_);
                }
            }
            CommandType::MultiCrate => {
                let mut temp: VecDeque<char> = VecDeque::new();
                for _ in 0..command.count {
                    match self.stacks[command.source].pop_front() {
                        Some(c) => match command_type {
                            CommandType::SingleCrate => self.stacks[command.target].push_front(c),
                            CommandType::MultiCrate => temp.push_front(c),
                        },
                        None => {
                            return Err(Error::DataError(format!(
                                "Invalid command {:?}: Popped from empty stack {}",
                                command, command.source
                            )))
                        }
                    };
                }
                for crate_ in temp {
                    self.stacks[command.target].push_front(crate_);
                }
            }
        }

        Ok(())
    }
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

/// Enum representing in which way crates are moved.
#[derive(Debug)]
pub enum CommandType {
    /// Single create at a time is moved
    SingleCrate,
    /// All crates moved as one stack
    MultiCrate,
}

#[derive(Debug)]
pub struct Command {
    /// Index of source stack from which to move crates
    pub source: usize,
    /// Index of target stack onto which to move crates
    pub target: usize,
    /// Number of creates to move
    pub count: usize,
}

impl FromStr for Command {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        // This'll be recompiled on every call, which is not perfect, but good enough
        let re = Regex::new("^move ([0-9]+) from ([1-9][0-9]*) to ([1-9][0-9]*)$").unwrap();

        let matches = match re.captures(s) {
            Some(c) => c,
            None => return Err(Error::DataError(format!("Invalid command: {}", s))),
        };

        // Theoretically we should be able to unwrap those, as the regex above should ensure they
        // are valid.
        let count: usize = matches[1]
            .parse()
            .map_err(|e| Error::DataError(format!("Invalid command: {}", e)))?;

        let source: usize = matches[2]
            .parse()
            .map_err(|e| Error::DataError(format!("Invalid command: {}", e)))?;

        let target: usize = matches[3]
            .parse()
            .map_err(|e| Error::DataError(format!("Invalid command: {}", e)))?;

        // In commands, both are 1-indexed. We'll transform them into 0-indexed here.
        Ok(Command {
            count,
            source: source - 1,
            target: target - 1,
        })
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

    #[test]
    fn command_from_str() {
        let cmd = Command::from_str("move 27 from 4 to 42").unwrap();
        assert_eq!(27, cmd.count);
        assert_eq!(3, cmd.source);
        assert_eq!(41, cmd.target);
    }

    #[test]
    fn apply_command_single() -> Result<(), Error> {
        let mut supplies = Supplies {
            stacks: vec![
                VecDeque::from(vec!['N', 'Z']),
                VecDeque::from(vec!['D', 'C', 'M']),
                VecDeque::from(vec!['P']),
            ],
        };

        let cmd = Command {
            count: 2,
            source: 1,
            target: 0,
        };
        supplies.apply(&cmd, CommandType::SingleCrate)?;

        assert_eq!(
            supplies,
            Supplies {
                stacks: vec![
                    VecDeque::from(vec!['C', 'D', 'N', 'Z']),
                    VecDeque::from(vec!['M']),
                    VecDeque::from(vec!['P']),
                ],
            }
        );

        Ok(())
    }

    #[test]
    fn apply_command_multi() -> Result<(), Error> {
        let mut supplies = Supplies {
            stacks: vec![
                VecDeque::from(vec!['N', 'Z']),
                VecDeque::from(vec!['D', 'C', 'M']),
                VecDeque::from(vec!['P']),
            ],
        };

        let cmd = Command {
            count: 2,
            source: 1,
            target: 0,
        };
        supplies.apply(&cmd, CommandType::MultiCrate)?;

        assert_eq!(
            supplies,
            Supplies {
                stacks: vec![
                    VecDeque::from(vec!['D', 'C', 'N', 'Z']),
                    VecDeque::from(vec!['M']),
                    VecDeque::from(vec!['P']),
                ],
            }
        );

        Ok(())
    }
}
