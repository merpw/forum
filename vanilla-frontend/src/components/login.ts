import { LoginForm, SignupForm } from "LoginModule"
import { Auth } from "./auth.js"

// All functionality for the login/signup form
export const loginController = () => {
  const loginSignupForm = document.querySelector(".container") as HTMLElement,
    pwShowHide = document.querySelectorAll(
      ".showHidePw"
    ) as NodeListOf<HTMLElement>,
    pwFields = document.querySelectorAll(
      ".password"
    ) as NodeListOf<HTMLInputElement>,
    signUpLink = document.querySelector(".signup-link") as HTMLInputElement,
    loginLink = document.querySelector(".login-link") as HTMLInputElement

  loginSignupForm.classList.remove("active")

  // To show/hide password in Auth form.
  pwShowHide.forEach((eyeIcon) => {
    eyeIcon.addEventListener("click", () => {
      pwFields.forEach((pwField) => {
        if (pwField.type === "password") {
          pwField.type = "text"
          pwShowHide.forEach((icon) => {
            icon.classList.replace("uil-eye-slash", "uil-eye")
          })
        } else {
          pwField.type = "password"
          pwShowHide.forEach((icon) => {
            icon.classList.replace("uil-eye", "uil-eye-slash")
          })
        }
      })
    })
  })
  // To go from login to sign in.
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
    // TODO: Error handling here maybe?
  }
}

class Login {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private onSubmit(event: Event) {
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
          Auth(true)
          // return response.json();
        } else {
          // login unsuccessful
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
    const usernameInput =
        this.form.querySelector<HTMLInputElement>("#username-register"),
      emailInput = this.form.querySelector<HTMLInputElement>("#email"),
      passwordInput =
        this.form.querySelector<HTMLInputElement>("#password-register"),
      firstNameInput = this.form.querySelector<HTMLInputElement>("#first-name"),
      lastNameInput = this.form.querySelector<HTMLInputElement>("#last-name"),
      ageInput = this.form.querySelector<HTMLInputElement>("#age"),
      genderInput = this.form.querySelector<HTMLInputElement>("#gender"),
      passwordRepeatInput = this.form.querySelector<HTMLInputElement>(
        "#password-register-repeat"
      )

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
