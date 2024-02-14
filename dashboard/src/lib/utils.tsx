import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { toast } from "sonner";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function DisplayAPIErrorToast(title: string, description: string, error: APIError) {
    toast.error(title, {
        description: <div className="flex flex-col gap-1"><p>{description}</p><p>Code: {error.code}, Message: {error.message}</p></div>,
        classNames: {
            toast: 'group-[.toaster]:bg-red-900',
        }
    });
}

export function DisplayExceptionErrorToast(title: string, description: string, exception: any) {
    toast.error(title, {
        description: <div className="flex flex-col gap-1"><p>{description}</p><p>{exception.message}</p></div>,
        classNames: {
            toast: 'group-[.toaster]:bg-red-900',
        }
    });
}