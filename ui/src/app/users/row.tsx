import { useState, Dispatch, SetStateAction, useEffect, useRef, useContext } from "react"
import { UseMutationResult, useMutation, useQueryClient } from "react-query"
import { RequiredCSRParams, deleteUser } from "../queries"
import { ConfirmationModalData, ConfirmationModal, ChangePasswordModalData, ChangePasswordModal } from "./components"
import "./../globals.scss"
import { useAuth } from "../auth/authContext"
import { AsideContext } from "../aside"

type rowProps = {
    id: number,
    username: string,

    ActionMenuExpanded: number
    setActionMenuExpanded: Dispatch<SetStateAction<number>>
}

export default function Row({ id, username, ActionMenuExpanded, setActionMenuExpanded }: rowProps) {
    const auth = useAuth()
    const asideContext = useContext(AsideContext)
    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData>(null)
    const [changePasswordModalData, setChangePasswordModalData] = useState<ChangePasswordModalData>(null)
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteUser, {
        onSuccess: () => queryClient.invalidateQueries('users')
    })
    const mutationFunc = (mutation: UseMutationResult<any, unknown, RequiredCSRParams, unknown>, params: RequiredCSRParams) => {
        mutation.mutate(params)
    }
    const handleDelete = () => {
        setConfirmationModalData({
            onMouseDownFunc: () => mutationFunc(deleteMutation, { id: id.toString(), authToken: auth.user ? auth.user.authToken : "" }),
            warningText: `Deleting user: "${username}". This action cannot be undone.`
        })
    }
    const handleChangePassword = () => {
        // asideContext.setExtraData({ "user": { "id": id, "username": username } })
        // asideContext.setIsOpen(true)
        setChangePasswordModalData({ "id": id.toString(), "username": username })
    }

    return (
        <>
            <tr>
                <td className="" width={5} aria-label="id">{id}</td>
                <td className="" aria-label="username">{username}</td>
                <td className="has-overflow" data-heading="Actions">
                    <span className="p-contextual-menu--center u-no-margin--bottom">
                        <button
                            className="p-contextual-menu__toggle p-button--base is-small u-no-margin--bottom"
                            aria-label="action-menu-button"
                            aria-controls="action-menu"
                            aria-expanded={ActionMenuExpanded == id ? "true" : "false"}
                            aria-haspopup="true"
                            onClick={() => setActionMenuExpanded(id)}
                            onBlur={() => setActionMenuExpanded(0)}>
                            <i className="p-icon--menu p-contextual-menu__indicator"></i>
                        </button>
                        <span className="p-contextual-menu__dropdown" id="action-menu" aria-hidden={ActionMenuExpanded == id ? "false" : "true"}>
                            <span className="p-contextual-menu__group">
                                {id == 1 ?
                                    <button className="p-contextual-menu__link" onMouseDown={handleDelete} disabled={true}>Delete User</button>
                                    :
                                    <button className="p-contextual-menu__link" onMouseDown={handleDelete}>Delete User</button>
                                }
                                <button className="p-contextual-menu__link" onMouseDown={handleChangePassword} >Change Password</button>
                            </span>
                        </span>
                    </span>
                </td>
                {confirmationModalData != null && <ConfirmationModal modalData={confirmationModalData} setModalData={setConfirmationModalData} />}
                {changePasswordModalData != null && <ChangePasswordModal modalData={changePasswordModalData} setModalData={setChangePasswordModalData} />}
            </tr>
        </>
    )
}