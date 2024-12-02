let filePath = "./input.txt";

if (Deno.args) {
  if (Deno.args[0] === "test") {
    console.log("Loading sample input.");
    filePath = "./sample.txt";
  }
}

let safeReportsPart1 = 0;
let safeReportsPart2 = 0;

function check(items: number[]): [boolean, number] {
  // If all goes in same direction, this should be +5 or -5.
  // const direction = (items[1] - items[0]) / Math.abs(items[1] - items[0]);
  let direction = 0;
  for (let i = 1; i < items.length; i++) {
    const diff = items[i] - items[i - 1];
    if (Math.abs(diff) > 3) {
      // Unsafe!
      return [false, i];
    }
    if (diff === 0) {
      // Unsafe!
      return [false, i];
    }

    if (diff > 0) {
      if (direction < 0) {
        // Unsafe!
        return [false, i];
      }
      direction = 1;
    }

    if (diff < 0) {
      if (direction > 0) {
        // Unsafe!
        return [false, i];
      }
      direction = -1;
    }
  }
  return [true, 0];
}

try {
  const data = await Deno.readTextFile(filePath);
  data.split("\n").forEach((line) => {
    // Trim and skip empty lines
    const trimmedLine = line.trim();
    if (!trimmedLine) return;

    const items = line.split(" ").map((item) => Number(item));
    let [isSafe, unsafeLevel] = check(items);
    if (isSafe) {
      safeReportsPart1++;
      safeReportsPart2++;
      console.log(`${trimmedLine} is safe`);
      return;
    }

    // Not safe.
    console.log(`${trimmedLine}: level ${unsafeLevel} makes it unsafe.`);

    for (let i = 0; i < items.length; i++) {
      const newLine = items.toSpliced(i, 1);
      [isSafe] = check(newLine);
      if (!isSafe) {
        console.log(`\tRemoving ${i + 1} (${newLine.join(',')}) does not make it safe.`);
        continue;
      }
      safeReportsPart2++;
      console.log(`\t${trimmedLine} is safe after removing item ${i + 1}`);
      return;
    }
    console.log(`\t${trimmedLine} is never safe`);
  });
} catch (error) {
  console.error("Error reading file:", error);
}

console.log("Part 1:", safeReportsPart1);
console.log("Part 1:", safeReportsPart2);
