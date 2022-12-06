use std::str::FromStr;

use crate::error::Error;

pub struct Range {
    pub start: usize,
    pub end: usize,
}

impl Range {
    /// Return whether this range fully contains the other.
    ///
    /// This relation is not reflexive. If a contains b, then b will usually not contain a - the
    /// sole exception being a = b.
    pub fn contains(&self, other: &Range) -> bool {
        self.start <= other.start && self.end >= other.end
    }

    /// Return whether this range overlaps the other.
    ///
    /// If the two share a border then this is treated as an overlap. As such these should not be
    /// thought of as ranges over the reals, but rather ranges over some discrete domain.
    ///
    /// This relation is reflexive. If a overlaps b, b will also overlap a.
    pub fn overlaps(&self, other: &Range) -> bool {
        let larger_start = self.start.max(other.start);
        let smaller_end = self.end.min(other.end);

        // Two ranges overlap if the interval between the larger of their starts, and smaller of
        // their ends, is not empty.
        // As we treat them sharing a border as overlapping, it's a >= rather than the >.
        smaller_end >= larger_start
    }
}

impl FromStr for Range {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let parts: Vec<&str> = s.split("-").collect();
        if parts.len() != 2 {
            return Err(Error::DataError(format!(
                "Invalid input, expected two numbers separated by hyphen, got: {}",
                s
            )));
        }

        let start: usize = parts[0].parse().map_err(|_| {
            Error::DataError(format!("Invalid input for start of range: {}", parts[0]))
        })?;
        let end: usize = parts[1].parse().map_err(|_| {
            Error::DataError(format!("Invalid input for end of range: {}", parts[1]))
        })?;

        Ok(Range { start, end })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_contains() {
        let r1 = Range { start: 1, end: 8 };

        // Case 1: Disjoint
        let r2 = Range { start: 9, end: 16 };
        assert!(!r1.contains(&r2));
        assert!(!r2.contains(&r1));

        // Case 2a: Overlapping but not contained
        let r2 = Range { start: 2, end: 9 };
        assert!(!r1.contains(&r2));
        assert!(!r2.contains(&r1));

        // Case 2b: Overlapping but not contained, the other way
        let r2 = Range { start: 0, end: 7 };
        assert!(!r1.contains(&r2));
        assert!(!r2.contains(&r1));

        // Case 2c: Touching at the very border
        let r2 = Range { start: 8, end: 9 };
        assert!(!r1.contains(&r2));
        assert!(!r2.contains(&r1));

        // Case 3: Fully contained
        let r2 = Range { start: 2, end: 8 };
        assert!(r1.contains(&r2));
        assert!(!r2.contains(&r1));

        // Case 4: Ranges are equal
        let r2 = Range { start: 1, end: 8 };
        assert!(r1.contains(&r2));
        assert!(r2.contains(&r1));
    }

    #[test]
    fn test_overlaps() {
        let r1 = Range { start: 1, end: 8 };

        // Case 1: Disjoint
        let r2 = Range { start: 9, end: 16 };
        assert!(!r1.overlaps(&r2));
        assert!(!r2.overlaps(&r1));

        // Case 2a: Overlapping but not contained
        let r2 = Range { start: 2, end: 9 };
        assert!(r1.overlaps(&r2));
        assert!(r2.overlaps(&r1));

        // Case 2b: Overlapping but not contained, the other way
        let r2 = Range { start: 0, end: 7 };
        assert!(r1.overlaps(&r2));
        assert!(r2.overlaps(&r1));

        // Case 2c: Touching at the very border
        let r2 = Range { start: 8, end: 9 };
        assert!(r1.overlaps(&r2));
        assert!(r2.overlaps(&r1));

        // Case 3: Fully contained
        let r2 = Range { start: 2, end: 8 };
        assert!(r1.overlaps(&r2));
        assert!(r2.overlaps(&r1));
    }
}
