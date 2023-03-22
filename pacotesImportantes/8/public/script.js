var startbutton = document.getElementById("start");
var wthere = document.getElementById("wthere");
var spanid = document.getElementById("spanid");
var who = document.getElementById("who");

who.style.visibility = "hidden";
wthere.style.visibility = "hidden";

startbutton.addEventListener("click", (e) => {
  startbutton.style.visibility = "hidden";
  spanid.textContent = "Knock Knock";
  wthere.style.visibility = "visible";
});

who.addEventListener("click", (e) => {
  who.style.visibility = "hidden";
  spanid.textContent = wo[randomm];
});

// who.textContent = wothere[randomm] + " who?";

wthere.addEventListener("click", (e) => {
  wthere.style.visibility = "hidden";
  spanid.textContent = "Knock Knock";
  who.style.visibility = "visible";
  spanid.textContent = wothere[randomm];
});

var wothere = [
  "You",
  "Butter",
  "Tank",
  "Annie",
  "Spell",
  "A little old lady",
  "Woo",
  "Sadie",
  "Ash"
];

var wo = [
  "You-hoo! Anybody home?",
  "Butter let me in or I'll freeze",
  "You're welcome",
  "Annie way you can let me in?",
  "W.H.O.",
  "Hey, you can yodel!",
  "Glad you're exited too!",
  "Sadie magic words and watch me dissapear",
  "Sounds like you have a cold!"
];

var randomm = Math.floor(Math.random() * 8);
who.textContent = wothere[randomm] + " who?";
