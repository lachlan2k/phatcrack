import { AxiosError } from 'axios'
import { useToast } from "vue-toastification"

export function useToastError() {
    const toast = useToast()

    const catcher = (e: any) => {
        let errorString = 'Unknown Error'
        if (e instanceof AxiosError) {
        errorString = e.response?.data?.message
        } else if (e instanceof Error) {
        errorString = e.message
        }

        toast.error(errorString)
    }

    return {
        catcher
    }
}