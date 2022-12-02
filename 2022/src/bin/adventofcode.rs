use std::{path::PathBuf, process::exit, str::FromStr};

use adventofcode::{calories, elves::Elf, error::Error, input, rps::Round};
use log::{error, info};

fn main() {
    env_logger::init();

    info!("Advent of Code: 2022");

    // match day_one() {
    //     Ok(_) => info!(""),
    //     Err(e) => abort(e),
    // }

    match day_two() {
        Ok(_) => info!(""),
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

    let mut score = 0;
    let rounds = input::to_string(input::input_path(2))?;
    let rounds: Vec<&str> = rounds.lines().collect();

    for round in rounds {
        let round = Round::from_str(round)?;
        // We are player two
        let (_, score_b) = round.score();
        score += score_b;
    }

    info!("Final score if following strategy book: {}", score);

    Ok(())
}

fn abort(e: Error) {
    error!("Error: {}", e);
    exit(1);
}
