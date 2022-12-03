const text = document.querySelector("pre").textContent;
const lines = text.trim().split('\n');

const scoreMapA = {
  "A X": 4, // 3 + 1 // 1 1
  "A Y": 8, // 6 + 2 // 1 2
  "A Z": 3, // 0 + 3 // 1 3
  "B X": 1, // 0 + 1 // 2 1
  "B Y": 5, // 3 + 2 // 2 2
  "B Z": 9, // 6 + 3 // 2 3
  "C X": 7, // 6 + 1 // 3 1
  "C Y": 2, // 0 + 2 // 3 2
  "C Z": 6, // 3 + 3 // 3 3
};
const totalScore = lines.map(l => scoreMapA[l]).reduce((partialSum, a) => partialSum + a, 0);
console.log(totalScore);