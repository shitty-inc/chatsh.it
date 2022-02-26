import { useEffect, createRef } from 'react';
import Clipboard from 'clipboard';

interface LinkProps {
  id: string | undefined;
  display: boolean;
}

function Link(props: LinkProps) {
  const downloadRef: React.RefObject<HTMLInputElement> = createRef();
  const copyRef: React.RefObject<HTMLButtonElement> = createRef();

  const { id, display } = props;

  useEffect(() => {
    if(display && id) {
      const downloadLink: HTMLInputElement = downloadRef.current as HTMLInputElement;
      const copyButton: HTMLButtonElement = copyRef.current as HTMLButtonElement;

      const clipboard = new Clipboard(
        copyButton, {
          target: () => downloadLink
        }
      )

      downloadLink.select();

      return () => clipboard.destroy();
    }
  });

  if(!display || !id) {
    return null;
  }

  return (
    <div className="row">
      <div className="col-6 offset-md-3">
        <p className="text-center">Copy this shit</p>
        <div className="input-group input-group-sm">
          <input
            readOnly={ true }
            type="text"
            ref={ downloadRef }
            className="form-control"
            value={ `${window.location.origin}/#/${id}` }
          />
          <button
            className="btn btn-outline-secondary"
            type="button"
            ref={ copyRef }
          >
            Copy
          </button>
        </div>
      </div>
    </div>
  );
}

export default Link;
