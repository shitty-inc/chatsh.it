interface FooterProps {
  status: string;
}

function Footer(props: FooterProps) {
  const { status } = props;

  return (
    <footer id="footer">
      <div className="container text-center">
        <p>{status}</p>
      </div>
    </footer>
  );
}

export default Footer;
