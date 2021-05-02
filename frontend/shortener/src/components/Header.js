import PropTypes from 'prop-types'


const Header = ({title, onShorten}) => {
    
    return (
        <header className="header">
            <h1>{title}</h1>            
            <button className="btn" onClick={onShorten}>Shorten!</button>
        </header>
    )
}

Header.defaultProps = {
    title: "URL Shortener",
}

Header.protoTypes = {
    title: PropTypes.string.isRequired,
}

export default Header
