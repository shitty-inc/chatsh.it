interface MessagesProps {
  outgoingText: string;
  setOutgoingtext: (text: string) => void;
  handleSubmit: (event: React.SyntheticEvent) => void;
  display: boolean;
}

function Messages(props: MessagesProps) {
  const { outgoingText, setOutgoingtext, handleSubmit } = props;

  if(!props.display) {
    return null;
  }

  return (
    <div className="row">
      <div className="col"></div>
      <div className="col-5">
        <form className="form" onSubmit={ handleSubmit }>
          <div className="input-group input-group-sm mb-3">
            <textarea
              className="form-control outgoing"
              value = { outgoingText }
              onChange={ e => setOutgoingtext(e.target.value) }>
            </textarea>
            <button className="btn btn-outline-secondary" type="submit" id="button-addon2">Send</button>
          </div>
        </form>
      </div>
      <div className="col"></div>
    </div>
  );
}

export default Messages;
