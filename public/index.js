"use strict";


const darkModeSwitch = document.querySelector("#darkmode-switch");

let colorScheme =
    window.matchMedia("(prefers-color-scheme: dark)").matches?
        "dark" :
        "light";


const storedMode = localStorage.getItem("darkMode");
if (storedMode) {
    colorScheme = storedMode;
}

switchColorScheme(colorScheme);

darkModeSwitch.addEventListener("change", (event) => {
    if(darkModeSwitch.checked)
        return switchColorScheme("dark");

    return switchColorScheme("light");
})



/**
 * switches the color scheme to either dark or ligt
 * @param {"light" | "dark"} mode
 */
function switchColorScheme(mode) {
    if(mode !== "dark" && mode !== "light") {
        console.log("unknown color scheme: ", mode);
        console.log(mode === "light")
        return;
    }

    console.log("switching color scheme: ", mode);

    if(mode === "dark") {
        document.body.classList.remove("light");
        document.body.classList.add("dark");
        darkModeSwitch.checked = true;
        colorScheme = "dark";
        localStorage.setItem("darkMode", "dark");
    } else {
        document.body.classList.remove("dark");
        document.body.classList.add("light");
        colorScheme = "light";
        localStorage.setItem("darkMode", "light");
    }
}