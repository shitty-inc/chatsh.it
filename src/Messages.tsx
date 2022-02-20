import { Message } from './App';

interface MessagesProps {
  messages: Message[];
}

function Messages(props: MessagesProps) {
  const { messages } = props;

  if(messages.length === 0) {
    return null;
  }

  return (
    <div className="row">
      <div className="col"></div>
      <div className="col-5">
        <div className="outgoing">
          <div className="chat card">
            <div className="card-body height3">
              <ul className="chat-list">
                { messages.slice(-5).map((message: Message, index: number) => {
                  return(<li className={ message.direction } key={ index }>
                    <div className="chat-body">
                      <div className="chat-message">
                        <h5>{ message.timestamp }</h5>
                        <p>{ message.text }</p>
                      </div>
                    </div>
                  </li>);
                }) }
              </ul>
            </div>
          </div>
        </div>
      </div>
      <div className="col"></div>
    </div>
  );
}

export default Messages;
