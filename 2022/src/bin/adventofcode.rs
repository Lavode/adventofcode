use std::{collections::HashSet, process::exit, str::FromStr};

use adventofcode::{
    calories, elves::Elf, error::Error, input, range::Range, rps::Round, rucksack::Rucksack,
};
use log::{debug, error, info};

fn main() {
    env_logger::init();

    info!("Advent of Code: 2022");

    let day: usize = std::env::args()
        .nth(1)
        .expect("Usage: adventofcode <day>")
        .parse()
        .expect("Invalid value for parameter <day>, must be numerical");

    let day_fn = match day {
        1 => day_one,
        2 => day_two,
        3 => day_three,
        4 => day_four,
        _ => panic!("Invalid day selected"),
    };

    match day_fn() {
        Ok(_) => {}
        Err(e) => abort(e),
    }
}

fn day_one() -> Result<(), Error> {
    info!("Day 1");
    let calories = calories::from_file(input::input_path(1))?;

    let mut elves: Vec<Elf> = calories
        .iter()
        .map(|cal| Elf {
            calories: cal.clone(),
        })
        .collect();

    elves.sort_by(|a, b| a.total_calories().cmp(&b.total_calories()));
    // sort_by sorts in ascendin order
    elves.reverse();

    let loaded_elf = &elves[0];

    info!(
        "Elf with most calories has total of {} cal in inventory",
        loaded_elf.total_calories()
    );

    let top_three = &elves[0..3];
    let top_three_cals: u32 = top_three.iter().map(|e| e.total_calories()).sum();
    info!(
        "Top three elves with most calories have total of {} cal in inventory",
        top_three_cals,
    );

    Ok(())
}

fn day_two() -> Result<(), Error> {
    info!("Day 2");

    // Score if using rock-paper-scissors serialization
    let mut score_rps = 0;
    // Score if using lose-draw-win serialization
    let mut score_ldw = 0;

    let rounds = input::to_string(input::input_path(2))?;
    let rounds: Vec<&str> = rounds.lines().collect();

    for round in rounds {
        let round_rps = Round::parse(
            round,
            adventofcode::rps::RoundSerializationFormat::RockPaperScissors,
        )?;

        let round_ldw = Round::parse(
            round,
            adventofcode::rps::RoundSerializationFormat::LoseDrawWin,
        )?;

        // We are player two
        let (_, score_b_rps) = round_rps.score();
        let (_, score_b_ldw) = round_ldw.score();

        score_rps += score_b_rps;
        score_ldw += score_b_ldw;
    }

    info!(
        "Final score if following strategy book with rock-paper-scissor format: {}",
        score_rps
    );
    info!(
        "Final score if following strategy book with lose-draw-win format: {}",
        score_ldw
    );

    Ok(())
}

fn day_three() -> Result<(), Error> {
    info!("Day 3");
    day_three_rucksack_packing()?;
    day_three_badge_finding()?;

    Ok(())
}

fn day_three_rucksack_packing() -> Result<(), Error> {
    let rucksacks = input::to_string(input::input_path(3))?;
    let rucksacks: Vec<&str> = rucksacks.lines().collect();

    let mut sum = 0;

    for items in rucksacks {
        let rucksack = Rucksack::from_str(items)?;
        let intersection = rucksack.intersection();
        if intersection.len() != 1 {
            return Err(Error::DataError(format!(
                "Expected one item to be in both compartments, got {}",
                intersection.len()
            )));
        }

        let item = intersection[0];
        if !item.is_ascii_alphabetic() {
            return Err(Error::DataError(format!(
                "Expected item to be ASCII alphabetic, was not: {}",
                item
            )));
        }

        let priority = item_priority(item);

        debug!(
            "Common item in both backpacks: {}, has priority: {}",
            item, priority
        );

        sum += priority;
    }

    info!("Sum of items weighed by priority: {}", sum);

    Ok(())
}

fn day_three_badge_finding() -> Result<(), Error> {
    let rucksacks = input::to_string(input::input_path(3))?;
    let rucksacks: Vec<&str> = rucksacks.lines().collect();

    if rucksacks.len() % 3 != 0 {
        return Err(Error::DataError(format!(
            "Number of rucksacks must be multiple of three; was: {}",
            rucksacks.len()
        )));
    }

    let mut sum = 0;

    for idx in 0..rucksacks.len() / 3 {
        let items_a = Rucksack::from_str(rucksacks[idx * 3 + 0])?.items();
        let items_b = Rucksack::from_str(rucksacks[idx * 3 + 1])?.items();
        let items_c = Rucksack::from_str(rucksacks[idx * 3 + 2])?.items();

        let badge_items: HashSet<char> = items_a.intersection(&items_b).map(|c| *c).collect();
        let badge_items: HashSet<char> = badge_items.intersection(&items_c).map(|c| *c).collect();

        if badge_items.len() != 1 {
            return Err(Error::DataError(format!(
                "Expected exactly one item to be shared between all three elves, got {}",
                badge_items.len()
            )));
        }

        let badge: Vec<&char> = badge_items.iter().take(1).collect();
        let badge = badge[0];
        let priority = item_priority(*badge);

        debug!(
            "Badge item determined to be {}, with priority {}",
            badge, priority
        );

        sum += priority;
    }

    info!("Sum of badge priorities: {}", sum);

    Ok(())
}

fn item_priority(item: char) -> u32 {
    let ascii_lowercase_a = 'a' as u32;
    let ascii_uppercase_a = 'A' as u32;

    // a-z => 1-26
    // A-Z => 27-52
    return match item.is_ascii_lowercase() {
        true => item as u32 - ascii_lowercase_a + 1,
        false => item as u32 - ascii_uppercase_a + 27,
    };
}

fn day_four() -> Result<(), Error> {
    let ranges_list = input::to_string(input::input_path(4))?;
    let ranges_list: Vec<&str> = ranges_list.lines().collect();

    let mut contained_count = 0;
    let mut overlapping_count = 0;

    for range_tuple in ranges_list {
        let ranges: Vec<&str> = range_tuple.split(",").collect();
        if ranges.len() != 2 {
            return Err(Error::DataError(format!(
                "Expected line to contain two ranges, got: {}",
                range_tuple
            )));
        }

        let r1 = Range::from_str(ranges[0])?;
        let r2 = Range::from_str(ranges[1])?;

        if r1.contains(&r2) || r2.contains(&r1) {
            contained_count += 1;
        }

        if r1.overlaps(&r2) {
            overlapping_count += 1;
        }
    }

    info!(
        "For {} pairs of ranges, one contains the other",
        contained_count
    );

    info!("For {} pairs of ranges, they overlap", overlapping_count);

    Ok(())
}

fn abort(e: Error) {
    error!("Error: {}", e);
    exit(1);
}
