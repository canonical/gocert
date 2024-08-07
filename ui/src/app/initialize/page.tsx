"use client"

import { useState, ChangeEvent } from "react"
import { getStatus, login, postFirstUser } from "../queries"
import { useMutation, useQuery } from "react-query"
import { useRouter } from "next/navigation"
import { passwordIsValid } from "../utils"
import { useAuth } from "../auth/authContext"
import { Logo } from "../nav"
import { useCookies } from "react-cookie"
import { statusResponse } from "../types"


export default function Initialize() {
    const router = useRouter()
    const auth = useAuth()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const statusQuery = useQuery<statusResponse, Error>({
        queryFn: () => getStatus()
    })
    if (statusQuery.data && statusQuery.data.initialized) {
        auth.setFirstUserCreated(true)
        router.push("/login")
    }
    const loginMutation = useMutation(login, {
        onSuccess: (e) => {
            setErrorText("")
            setCookie('user_token', e, {
                sameSite: true,
                secure: true,
                expires: new Date(new Date().getTime() + 60 * 60 * 1000),
            })
            router.push('/certificate_requests')
        },
        onError: (e: Error) => {
            setErrorText(e.message)
        }
    })
    const postUserMutation = useMutation(postFirstUser, {
        onSuccess: () => {
            setErrorText("")
            auth.setFirstUserCreated(true)
            loginMutation.mutate({ username: username, password: password1 })
        },
        onError: (e: Error) => {
            setErrorText(e.message)
        }
    })
    const [username, setUsername] = useState<string>("")
    const [password1, setPassword1] = useState<string>("")
    const [password2, setPassword2] = useState<string>("")
    const passwordsMatch = password1 === password2
    const [errorText, setErrorText] = useState<string>("")
    const [showPassword1, setShowPassword1] = useState<boolean>(false)
    const handleUsernameChange = (event: ChangeEvent<HTMLInputElement>) => { setUsername(event.target.value) }
    const handlePassword1Change = (event: ChangeEvent<HTMLInputElement>) => { setPassword1(event.target.value) }
    const handlePassword2Change = (event: ChangeEvent<HTMLInputElement>) => { setPassword2(event.target.value) }
    return (
        <>
            <div style={{ backgroundColor: "#262626", height: "6.5vh" }}>
                <div style={{ marginLeft: "30px" }}>
                    <Logo />
                </div>
            </div>
            <div style={{
                display: "flex",
                alignContent: "center",
                justifyContent: "center",
                flexWrap: "wrap",
                height: "93.5vh",
            }}>
                <div className="p-panel" style={{
                    width: "45rem",
                    minWidth: "min-content",
                    minHeight: "min-content",
                }}>
                    <fieldset>
                        <div className="p-panel__header">
                            <h2>Initialize GoCert</h2>
                        </div>
                        <div className="p-panel__content">
                            <form className={"p-form-validation " + ((!passwordIsValid(password1) && password1 != "") || (!passwordsMatch && password2 != "") ? "is-error" : "")}>
                                <div className="p-form__group row">
                                    <h4>Create the initial admin user</h4>
                                    <label className="p-form__label">Username</label>
                                    <input type="text" id="InputUsername" name="InputUsername" onChange={handleUsernameChange} />
                                    <div>
                                        <label className="p-form__label">Password</label>
                                        <button className="p-button--base u-no-margin--bottom has-icon" style={{ float: "right" }} aria-live="polite" aria-controls="password" onClick={(e) => { e.preventDefault(); setShowPassword1(!showPassword1) }}>
                                            {showPassword1 ? (
                                                <>
                                                    <span className="p-form-password-toggle__label">
                                                        Hide
                                                    </span>
                                                    <i className="p-icon--hide"></i>
                                                </>
                                            ) : (
                                                <>
                                                    <span className="p-form-password-toggle__label">
                                                        Show
                                                    </span>
                                                    <i className="p-icon--show"></i>
                                                </>
                                            )}
                                        </button>
                                    </div>
                                    <input className="p-form-validation__input" type={showPassword1 ? "text" : "password"} id="password1" name="password" placeholder="******" autoComplete="current-password" required={true} onChange={handlePassword1Change} />
                                    <p className="p-form-help-text">
                                        Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol.
                                    </p>
                                    <label htmlFor="p-form__label">Confirm Password</label>
                                    <input className="p-form-validation__input" type="password" id="InputPassword" name="password2" placeholder="******" autoComplete="current-password" onChange={handlePassword2Change} />
                                    {!passwordIsValid(password1) && password1 != "" && <p className="p-form-validation__message">Password is not valid</p>}
                                    {passwordIsValid(password1) && !passwordsMatch && password2 != "" && <p className="p-form-validation__message">Passwords do not match</p>}
                                    {errorText &&
                                        <div className="p-notification--negative">
                                            <div className="p-notification__content">
                                                <h5 className="p-notification__title">Error</h5>
                                                <p className="p-notification__message">{errorText.split("error: ")}</p>
                                            </div>
                                        </div>
                                    }
                                    {!passwordsMatch || !passwordIsValid(password1) ? (
                                        <>
                                            <button className="p-button--positive" type="submit" name="submit" disabled={true}>Submit</button>
                                        </>
                                    ) : (
                                        <button className="p-button--positive" type="submit" name="submit" onClick={(event) => { event.preventDefault(); postUserMutation.mutate({ username: username, password: password1 }) }}>Submit</button>
                                    )}
                                </div>
                            </form>
                        </div>
                    </fieldset>
                </div >
            </div>
        </>
    )
}