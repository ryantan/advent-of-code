let filePath = "./input.txt";

if (Deno.args) {
  if (Deno.args[0] === "test") {
    console.log("Loading sample input.");
    filePath = "./sample.txt";
  }
}

const pattern = /mul\((\d{1,3}),(\d{1,3})\)/g;

let data = "";
try {
  data = await Deno.readTextFile(filePath);
  // console.log("data:", data);
} catch (error) {
  console.error("Error reading file:", error);
}

function getResult(data: string): number {
  let result = 0;
  let match;
  while ((match = pattern.exec(data)) !== null) {
    const firstNumber = parseInt(match[1], 10);
    const secondNumber = parseInt(match[2], 10);
    // console.log(`First number: ${firstNumber}, Second number: ${secondNumber}`);
    const product = firstNumber * secondNumber;
    result += product;
  }
  return result;
}

console.log("Part 1:", getResult(data));

try {
  if (Deno.args[0] === "test") {
    console.log("Loading sample input.");
    filePath = "./sample2.txt";
  }
  data = await Deno.readTextFile(filePath);
  // console.log("data:", data);
} catch (error) {
  console.error("Error reading file:", error);
}

const onlyDoParts: string[] = [];
for (let i = 0; i < 9999; i++) {
  const indexOfDoNot = data.indexOf("don't()");
  if (indexOfDoNot === -1) {
    // Didn't find any don't()? Add rest of data and finish.
    onlyDoParts.push(data);
    break;
  }
  onlyDoParts.push(data.substring(0, indexOfDoNot));
  data = data.substring(indexOfDoNot + 7);
  const indexOfDo = data.indexOf("do()");
  if (indexOfDo === -1) {
    // Didn't find any do()? Finish.
    break;
  }
  data = data.substring(indexOfDo + 4);
}
console.log("Part 2:", getResult(onlyDoParts.join(" ")));
