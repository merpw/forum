import { LoginForm, SignupForm } from "../types"
import { login, signup } from "../api/post.js"
import { Auth } from "./auth.js"


class Login {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private async onSubmit(event: Event) {
    event.preventDefault()
    const usernameInput = this.form.querySelector("#username-login") as HTMLInputElement,
    passwordInput = this.form.querySelector("#password-login") as HTMLInputElement,
    rememberMeInput = this.form.querySelector("#logCheck") as HTMLInputElement

    const formData: LoginForm = {
      login: usernameInput.value,
      password: passwordInput.value,
      rememberMe: rememberMeInput.checked,
    }
    const loginstatus = await login(formData)
    loginstatus == true ? Auth(true) : Auth(false) 
  }
}

class Signup {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private async onSubmit(event: Event) {
    event.preventDefault()
    // const formData: SignupForm = this.getFormData()
    await signup(this.getFormData() as SignupForm)
  }

  private getFormData(): SignupForm {
    const usernameInput = this.form.querySelector("#username-register") as HTMLInputElement,
    emailInput = this.form.querySelector("#email") as HTMLInputElement,
    passwordInput = this.form.querySelector("#password-register") as HTMLInputElement,
    firstNameInput = this.form.querySelector("#first-name") as HTMLInputElement,
    lastNameInput = this.form.querySelector("#last-name") as HTMLInputElement,
    ageInput = this.form.querySelector("#age") as HTMLInputElement,
    genderInput = this.form.querySelector("#gender") as HTMLInputElement,
    passwordRepeatInput = this.form.querySelector("#password-register-repeat") as HTMLInputElement

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
}

// All functionality for the login/signup form
export const loginController = async () => {
  const loginSignupForm = document.querySelector(".container") as HTMLElement,
  pwShowHide = document.querySelectorAll(".showHidePw") as NodeListOf<HTMLElement>,
  pwFields = document.querySelectorAll(".password") as NodeListOf<HTMLInputElement>,
  signUpLink = document.querySelector(".signup-link") as HTMLInputElement,
  loginLink = document.querySelector(".login-link") as HTMLInputElement

  if (loginSignupForm.classList.contains("active")) {
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

  const signupElement = document.querySelector(".signup") as HTMLDivElement

  signUpLink.addEventListener("click", () => {
    loginSignupForm.classList.add("active")
    signupElement.style.display = "block"
  })

  // To go from sign in to login.
  loginLink.addEventListener("click", () => {
    loginSignupForm.classList.remove("active")
    signupElement.style.display = "none"
  })

  new Login(document.querySelector("#login-form") as HTMLFormElement)
  new Signup(document.querySelector("#signup-form") as HTMLFormElement)
}
