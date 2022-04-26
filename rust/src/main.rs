use image::GenericImageView;
use std::fs::File;
use std::fs;
use std::env;
use std::process;
use std::io::{Write};

use curl::easy::Easy;

fn get_bcf(width: u32, height: u32) -> u32 {
  let mut bcf:u32 = 1;
  for i in 1..width {
    for j in 1..height {
      if width % i == 0 {
        if height % j == 0 && i == j {
          bcf = j;
        }
      }
    }
  }
  return bcf;
}

fn read_file<'a>(mut file: &'a str, r_flag: bool, req: bool, url: &'a mut String) -> i32 {
  let width: u32; let height: u32; let factor: u32;
  let wuneven: f32; let huneven: f32; 
  let img = image::open(file).unwrap();

  (width, height) = img.dimensions();

  factor = get_bcf(width, height);
  wuneven = (height as f32) / (width as f32);
  huneven = (width as f32) / (height as f32);

  if req == true {
    // split the url into a vector and then get the last element
    let mut url_vec: Vec<&str> = url.split("/").collect();
    file = url_vec.pop().unwrap();
  }

  if factor == 1 {
    if width < height {
      print!("{} > 1:{:.2} (uneven)", file, wuneven);
    } else {
      print!("{} > {:.2}:1 (uneven)", file, huneven);
    }
  } else {
    print!("{} > {}:{}", file, width/factor, height/factor);
  } 
  if r_flag == true {
    println!(" [{}x{}]", width, height);
  } else {
    println!("");
  }
  return 0;
}

fn download(url: &mut str) -> i32 {
  let mut handle = Easy::new();
  let mut file = File::create("/tmp/uirc.png").unwrap();
  handle.url(url).unwrap();
  handle.write_function(move |data| {
    file.write_all(data).unwrap();
    Ok(data.len())
  }).unwrap();
  handle.perform().unwrap();
  return 0;
}

fn handle_arg(arg: &mut String, mut r_flag: bool) -> bool {
  let complete:i32;
  let first:&str;
  let first_two:&str;
  let first_four:&str;
  let mut empty_str:String = String::from("");
  const VERSION:&str = "0.1.0";
  const HELP:&str = "USAGE: uirc [OPTIONS] IMAGE1 [IMAGE2] [...]

    OPTIONS:
        informational:
        -h, --help   \t: Display this message
        -l, --license\t: Display the license disclaimer foruirc-go (BSD 2-Clause)
        -v, --version\t: Display the version of uirc-go

        functional:
        -r, --res    \t: Display the resolution of the image (in addition to the ratio)

    HELP: 
        For more information on how to use uirc-go, open the man page uirc(1).";

  first = &arg[..1];

  if arg.chars().count() >= 4 {
    first_four = &arg[..4];
  } else {
    first_four = ""; 
  }

  if arg.chars().count() >= 2 {
    first_two = &arg[..2];
  } else {
    first_two = "";
  }

  if "--" == first_two || "-" == first {
    if "--help" == arg || "-h" == arg {
      println!("an unneccessary image ratio calculator, in Rust! (uirc-rust) v{}\n", VERSION);
      println!("Copyright 2022 Bryson Steck\nFree and Open Source under the BSD 2-Clause License\n");
      println!("{}", HELP);
      process::exit(1);
    } else if "--license" == arg || "-l" == arg {
      println!("uirc-rust is Free and Open Source Software under the BSD 2-Clause License.");
      println!("Please read the license regarding copying and distributing uirc.");
      println!("https://github.com/brysonsteck/uirc/blob/master/LICENSE");
      process::exit(1);
    } else if "--res" == arg || "-r" == arg {
      r_flag = true;
      return r_flag;
    } else if "--version" == arg || "-v" == arg {
      println!("uirc-rust v{}", VERSION);
      process::exit(1);
    } else {
      println!("uirc-rust: invalid argument \"{}\"\nType \"uirc --help\" for help with arguments.", arg);
      process::exit(5);
    }
  }

  if "http" == first_four {
    print!("downloading \"{}\"...", arg);
    std::io::stdout().flush().unwrap();
    download(arg);
    println!("ok");
    complete = read_file("/tmp/uirc.png", r_flag, true, arg);
    if complete != 0 {
      read_file(arg, r_flag, false, &mut empty_str);
    } else {
      let result = fs::remove_file("/tmp/uirc.png");
      if result.is_err() {
        println!("failed to remove temporary file");
      }
    }
  } else {
    // if no more flags, run ratio calculations
    read_file(arg, r_flag, false, &mut empty_str);
  }

  return r_flag;
}

fn main() {
  let mut runs:i32 = 0;
  let init_args: Vec<String> = env::args().collect();
  let args = &init_args[1..];
  let mut r_flag:bool = false;

  if args.len() < 1 {
    println!("uirc-rust: at least one argument is required");
    process::exit(1);
  }

  for arg in args.iter() {
    r_flag = handle_arg(&mut String::from(arg), r_flag);
    runs = runs + 1;
  }

  if runs < 2 && r_flag == true {
    println!("uirc-rust: at least one file/url is required");
    process::exit(1);
  }
}
