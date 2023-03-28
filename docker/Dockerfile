FROM golang:1.19-alpine

ENV RECAPTCHA_KEY=6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI
ENV RECAPTCHA_SECRET=6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379
ENV DEFAULT_EXPIRATION=3600
ENV MAXIMUM_EXPIRATION=3600
ENV CUSTOM_LOGO=https://raw.githubusercontent.com/clglavan/private-notes/add_logo/logo.png
ENV NOTE_MAX_LENGTH_SERVER=20000
ENV NOTE_MAX_LENGTH_CLIENT=10000
ENV LANG_INDEX_TITLE="Private notes"
ENV LANG_INDEX_SUBTITLE="Send self-destructing private notes securely"
ENV LANG_INDEX_NOTE_PLACEHOLDER="Notes will self-destruct after they are read..."
ENV LANG_INDEX_PASSWORD="Password protect your note"
ENV LANG_INDEX_PASSWORD_PLACEHOLDER="enter your password"
ENV LANG_INDEX_EXPIRATION="Custom expiration time,1-60,default 60 min"
ENV LANG_INDEX_SEND_BUTTON="Encrypt & Send"
ENV LANG_SUCCESS_TITLE="Thank you for using private notes"
ENV LANG_SUCCESS_SUBTITLE="This note will self-destruct after it will be read. Click on it to copy to clipboard and send this link to the other party."
ENV LANG_SUCCESS_TOOLTIP="Click on it to copy to clipboard"
ENV LANG_CONFIRM_SUBTITLE="Do you want to decrypt this message now? It's contents will be lost forever"
ENV LANG_CONFIRM_SHOW_BUTTON="Show & Destroy"
ENV LANG_RESULT_TITLE="Thank you for using private notes"
ENV LANG_RESULT_SUBTITLE="This note has been destroyed, below is the only copy."
ENV LANG_RESULT_PASSWORD="This message is password protected"
ENV LANG_RESULT_PASSWORD_PLACEHOLDER="Enter your password"
ENV LANG_RESULT_TOOLTIP="Click on it to copy to clipboard"
ENV LANG_ERROR_TITLE="Private notes"
ENV LANG_ERROR_SUBTITLE="Note does not exist"
ENV LANG_ERRORBAG_EMPTY="You can't send an empty note"
ENV LANG_ERRORBAG_PASSWORD_REQUIRED="Secret password checked but not provided"
ENV LANG_ERRORBAG_EXPIRATION_REQUIRED="Expiration time checked but not provided"
ENV LANG_ERRORBAG_NOTE_TOO_LONG="Secret note is too long"
ENV LANG_ERRORBAG_EXPIRATION_TOO_LONG="Expiration time is too long"

WORKDIR /private-notes

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /private-notes/

RUN go build -o /private-notes/privateNotes ./main
RUN chmod 755 /private-notes/privateNotes

EXPOSE 80

CMD ["/private-notes/privateNotes" ]
