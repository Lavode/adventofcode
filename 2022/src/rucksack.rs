use std::{collections::HashSet, str::FromStr};

use crate::error::Error;

pub struct Rucksack {
    pub compartment_a: Vec<char>,
    pub compartment_b: Vec<char>,
}

impl FromStr for Rucksack {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        if s.len() % 2 == 1 {
            return Err(Error::DataError(
                "Invalid input for Rucksack, must be of even length".to_string(),
            ));
        }

        Ok(Rucksack {
            compartment_a: s[..s.len() / 2].chars().collect(),
            compartment_b: s[s.len() / 2..].chars().collect(),
        })
    }
}

impl Rucksack {
    /// Set of items contained in either of the two compartments, ignoring duplicates.
    pub fn items(&self) -> HashSet<char> {
        let a: HashSet<&char> = HashSet::from_iter(self.compartment_a.iter());
        let b: HashSet<&char> = HashSet::from_iter(self.compartment_b.iter());

        a.union(&b).map(|c| **c).collect()
    }

    /// List of items contained in both compartments.
    pub fn intersection(&self) -> Vec<char> {
        let items_a: HashSet<&char> = HashSet::from_iter(self.compartment_a.iter());
        let items_b: HashSet<&char> = HashSet::from_iter(self.compartment_b.iter());

        // Love me some double-dereferencing
        items_a.intersection(&items_b).map(|c| **c).collect()
    }
}
