// Import Deno file system module
let filePath = "./input.txt";

if (Deno.args) {
  if (Deno.args[0] === "test") {
    console.log("Loading sample input.");
    filePath = "./sample.txt";
  }
}
// console.log('filePath:', filePath);

// Left numbers
const firstNumbers: number[] = [];
// Right numbers
const secondNumbers: number[] = [];

try {
  const data = await Deno.readTextFile(filePath);
  data.split("\n").forEach((line) => {
    // Trim and skip empty lines
    const trimmedLine = line.trim();
    if (!trimmedLine) return;

    const [first, second] = trimmedLine.split("   ");

    // Parse the numbers and add them to respective arrays
    firstNumbers.push(Number(first));
    secondNumbers.push(Number(second));
  });
} catch (error) {
  console.error("Error reading file:", error);
}

// console.log("First Numbers:", firstNumbers);
// console.log("Second Numbers:", secondNumbers);

firstNumbers.sort((a, b) => a - b);
secondNumbers.sort((a, b) => a - b);

// console.log("First Numbers sorted:", firstNumbers);
// console.log("Second Numbers sorted:", secondNumbers);

let differences = 0;
for (let i = 0; i < firstNumbers.length; i++) {
  differences += Math.abs(secondNumbers[i] - firstNumbers[i]);
}

console.log("Part 1:", differences);

// Part 2

const rightFrequency = new Map<number, number>();
for (const secondNumber of secondNumbers) {
  rightFrequency.set(
    secondNumber,
    (rightFrequency.get(secondNumber) || 0) + 1,
  );
}
// console.log('rightFrequency:', rightFrequency)

let similarity = 0;
for (const firstNumber of firstNumbers) {
  similarity += firstNumber * (rightFrequency.get(firstNumber) || 0);
}

console.log("Part 2:", similarity);
