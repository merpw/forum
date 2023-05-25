import { LoginForm, SignupForm } from "./types"
import { Auth } from "./auth.js"
import { wsHandler } from "./ws.js"

class Login {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private async onSubmit(event: Event) {
    event.preventDefault()
    const usernameInput =
    this.form.querySelector<HTMLInputElement>("#username-login"),
    passwordInput =
    this.form.querySelector<HTMLInputElement>("#password-login"),
    rememberMeInput = this.form.querySelector<HTMLInputElement>("#logCheck")

    if (usernameInput && passwordInput && rememberMeInput) {
      const formData: LoginForm = {
        login: usernameInput.value,
        password: passwordInput.value,
        rememberMe: rememberMeInput.checked,
      }

      fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      }).then((response) => {
          if (response.ok) {
            wsHandler() 
            Auth(true)
          } else {
            // TODO: Fix this error handling. It is super bad.
            response.text().then((error) => {
              console.log(`Error: ${error}`)
            })
          }
        })
    }
  }
}

class Signup {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private onSubmit(event: Event) {
    event.preventDefault()
    const formData: SignupForm = this.getFormData()
    console.log(formData)

    fetch("/api/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (response.ok) {
          loginController()
        } else {
          response.text().then((error) => {
            console.log(`Error: ${error}`)
            // TODO: Displaying error message to user.
          })
        }
      })
      .catch((error) => {
        console.error(error)
      })
  }

  private getFormData(): SignupForm {
    // Form inputs
    const 
    usernameInput = this.form.querySelector<HTMLInputElement>("#username-register"),
    emailInput = this.form.querySelector<HTMLInputElement>("#email"),
    passwordInput = this.form.querySelector<HTMLInputElement>("#password-register"),
    firstNameInput = this.form.querySelector<HTMLInputElement>("#first-name"),
    lastNameInput = this.form.querySelector<HTMLInputElement>("#last-name"),
    ageInput = this.form.querySelector<HTMLInputElement>("#age"),
    genderInput = this.form.querySelector<HTMLInputElement>("#gender"),
    passwordRepeatInput = this.form.querySelector<HTMLInputElement>("#password-register-repeat")

    if (
      firstNameInput &&
        lastNameInput &&
        usernameInput &&
        emailInput &&
        passwordInput &&
        ageInput &&
        genderInput &&
        passwordInput &&
        passwordRepeatInput
    ) {
      if (passwordInput.value != passwordRepeatInput.value) {
        //TODO: Display error message to user.
        console.log("password != repeat")
        throw new Error("Kek")
      }
      const formData: SignupForm = {
        name: usernameInput.value,
        email: emailInput.value,
        password: passwordInput.value,
        first_name: firstNameInput.value,
        last_name: lastNameInput.value,
        dob: ageInput.value,
        gender: genderInput.value,
      }
      return formData
    }
    throw new Error("Could not find form input fields.")
  }
}

// All functionality for the login/signup form
export const loginController = () => {
  const loginSignupForm = document.querySelector(".container") as HTMLElement,
  pwShowHide = document.querySelectorAll(".showHidePw") as NodeListOf<HTMLElement>,
  pwFields = document.querySelectorAll(".password") as NodeListOf<HTMLInputElement>,  
  signUpLink = document.querySelector(".signup-link") as HTMLInputElement,
  loginLink = document.querySelector(".login-link") as HTMLInputElement

  if (loginSignupForm.classList.contains("active")){
    loginSignupForm.classList.remove("active")
  }

  pwShowHide.forEach((eyeIcon) => {
    eyeIcon.addEventListener("click", () => {

      pwFields.forEach((input) => {
        if (input.getAttribute("type") === "password") {
          input.setAttribute("type", "text")
        } else {
          input.setAttribute("type", "password")
        }

      })
    })
  }) 

  signUpLink?.addEventListener("click", () => {
    loginSignupForm.classList.add("active")
  })

  // To go from sign in to login.
  loginLink?.addEventListener("click", () => {
    loginSignupForm.classList.remove("active")
  })

  const signupForm = document.querySelector<HTMLFormElement>("#signup-form")
  const loginForm = document.querySelector<HTMLFormElement>("#login-form")
  if (signupForm && loginForm) {
    new Signup(signupForm)
    new Login(loginForm)
  } else {
    console.error("Something went wrong.")
    // TODO: error handling here
  }
}
