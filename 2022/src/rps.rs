use std::str::FromStr;

use crate::error::Error;

/// Possible choices in single round of rock-paper-scissors.
#[derive(Debug, PartialEq, Eq, Clone)]
pub enum Choice {
    Rock,
    Paper,
    Scissors,
}

impl Choice {
    /// How much this choice is worth towards our score.
    pub fn score(&self) -> u32 {
        match self {
            Self::Rock => 1,
            Self::Paper => 2,
            Self::Scissors => 3,
        }
    }

    /// Calculate outcome of playing this choice versus another.
    pub fn versus(&self, other: &Choice) -> Outcome {
        if &self.wins_vs() == other {
            return Outcome::Win;
        } else if &self.loses_vs() == other {
            return Outcome::Loss;
        } else {
            return Outcome::Draw;
        }
    }

    /// The choice we win against.
    pub fn wins_vs(&self) -> Choice {
        match self {
            Self::Rock => Self::Scissors,
            Self::Paper => Self::Rock,
            Self::Scissors => Self::Paper,
        }
    }

    /// The choice we lose against.
    pub fn loses_vs(&self) -> Choice {
        match self {
            Choice::Rock => Self::Paper,
            Choice::Paper => Self::Scissors,
            Choice::Scissors => Self::Rock,
        }
    }
}

impl FromStr for Choice {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "A" | "X" => Ok(Self::Rock),
            "B" | "Y" => Ok(Self::Paper),
            "C" | "Z" => Ok(Self::Scissors),
            _ => Err(Error::DataError(format!(
                "Invalid value for player choice: {}",
                s
            ))),
        }
    }
}

/// Possible outcomes in a single round of rock-paper-scissors.
pub enum Outcome {
    Win,
    Draw,
    Loss,
}

impl Outcome {
    /// How much this outcome is worth towards our score.
    pub fn score(&self) -> u32 {
        match self {
            Self::Win => 6,
            Self::Draw => 3,
            Self::Loss => 0,
        }
    }
}

/// Round consisting of two players playing against each other.
#[derive(Debug, Eq, PartialEq)]
pub struct Round {
    /// Player A's choice of move.
    pub player_a_choice: Choice,
    /// Player B's choice of move.
    pub player_b_choice: Choice,
}

/// Possible formats in which the strategies are serialized.
pub enum RoundSerializationFormat {
    /// A and X stand for Rock, B and Y for Paper, C and Z for Scissors, for the first respectively
    /// the second player.
    RockPaperScissors,
    /// A stands for Rock, B for Paper, C for Scissors.
    /// X tells the second player to lose, Y to draw, Z to win.
    LoseDrawWin,
}

impl Round {
    pub fn parse(s: &str, format: RoundSerializationFormat) -> Result<Round, Error> {
        let parts: Vec<&str> = s.split(" ").collect();
        if parts.len() != 2 {
            return Err(Error::DataError(format!(
                "Invalid round, must contain exactly one choice for each player. Got: {}",
                s
            )));
        }

        let player_a_choice = Choice::from_str(parts[0])?;

        match format {
            RoundSerializationFormat::RockPaperScissors => {
                let player_b_choice = Choice::from_str(parts[1])?;

                return Ok(Round {
                    player_a_choice,
                    player_b_choice,
                });
            }
            RoundSerializationFormat::LoseDrawWin => {
                let player_b_choice = match parts[1] {
                    "X" => player_a_choice.wins_vs(),  // Make sure we lose
                    "Y" => player_a_choice.clone(),    // Make sure we draw
                    "Z" => player_a_choice.loses_vs(), // Make sure we win
                    _ => {
                        return Err(Error::DataError(format!(
                            "Invalid strategy for us: {}",
                            parts[1]
                        )))
                    }
                };

                return Ok(Round {
                    player_a_choice,
                    player_b_choice,
                });
            }
        }
    }

    /// Calculate scores achieved by player A and B respectively. Returns a tuple where first item
    /// is player A's score, second item player B's.
    pub fn score(&self) -> (u32, u32) {
        let (mut score_a, mut score_b): (u32, u32) = (0, 0);

        score_a += self.player_a_choice.score();
        score_a += self.player_a_choice.versus(&self.player_b_choice).score();

        score_b += self.player_b_choice.score();
        score_b += self.player_b_choice.versus(&self.player_a_choice).score();

        (score_a, score_b)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn strategy_score() {
        // Rock vs X
        let s = Round {
            player_a_choice: Choice::Rock,
            player_b_choice: Choice::Rock,
        };
        assert_eq!(s.score(), (1 + 3, 1 + 3));

        let s = Round {
            player_a_choice: Choice::Rock,
            player_b_choice: Choice::Paper,
        };
        assert_eq!(s.score(), (1 + 0, 2 + 6));

        let s = Round {
            player_a_choice: Choice::Rock,
            player_b_choice: Choice::Scissors,
        };
        assert_eq!(s.score(), (1 + 6, 3 + 0));

        // Paper vs X
        let s = Round {
            player_a_choice: Choice::Paper,
            player_b_choice: Choice::Rock,
        };
        assert_eq!(s.score(), (2 + 6, 1 + 0));

        let s = Round {
            player_a_choice: Choice::Paper,
            player_b_choice: Choice::Paper,
        };
        assert_eq!(s.score(), (2 + 3, 2 + 3));

        let s = Round {
            player_a_choice: Choice::Paper,
            player_b_choice: Choice::Scissors,
        };
        assert_eq!(s.score(), (2 + 0, 3 + 6));

        // Scissors vs X
        let s = Round {
            player_a_choice: Choice::Scissors,
            player_b_choice: Choice::Rock,
        };
        assert_eq!(s.score(), (3 + 0, 1 + 6));

        let s = Round {
            player_a_choice: Choice::Scissors,
            player_b_choice: Choice::Paper,
        };
        assert_eq!(s.score(), (3 + 6, 2 + 0));

        let s = Round {
            player_a_choice: Choice::Scissors,
            player_b_choice: Choice::Scissors,
        };
        assert_eq!(s.score(), (3 + 3, 3 + 3));
    }

    #[test]
    fn parse_round_success() {
        let s = Round::parse("A X", RoundSerializationFormat::LoseDrawWin).unwrap();
        assert_eq!(
            s,
            Round {
                player_a_choice: Choice::Rock,
                player_b_choice: Choice::Scissors,
            }
        );

        let s = Round::parse("C Y", RoundSerializationFormat::LoseDrawWin).unwrap();
        assert_eq!(
            s,
            Round {
                player_a_choice: Choice::Scissors,
                player_b_choice: Choice::Scissors,
            }
        );

        let s = Round::parse("A X", RoundSerializationFormat::RockPaperScissors).unwrap();
        assert_eq!(
            s,
            Round {
                player_a_choice: Choice::Rock,
                player_b_choice: Choice::Rock,
            }
        );
    }

    #[test]
    fn parse_round_invalid() {
        Round::parse("A X A", RoundSerializationFormat::LoseDrawWin).unwrap_err();
        Round::parse("A U", RoundSerializationFormat::LoseDrawWin).unwrap_err();
        Round::parse("A", RoundSerializationFormat::LoseDrawWin).unwrap_err();
    }
}
