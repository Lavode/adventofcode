#[derive(Clone, Debug)]
pub struct Elf {
    pub calories: Vec<u32>,
}

impl Elf {
    pub fn total_calories(&self) -> u32 {
        self.calories.iter().sum()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn total_calories_sums_calories() {
        let elf = Elf { calories: vec![] };
        assert_eq!(elf.total_calories(), 0);

        let elf = Elf {
            calories: vec![1, 42, 17],
        };
        assert_eq!(elf.total_calories(), 60);
    }
}
