package nntpclient

import (
	"errors"
	"io"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/EdasL/NNTP/nntp"
)

// Client is an NNTP client.
type Client struct {
	//textual network protocol connection
	conn   *textproto.Conn
	Banner string
}

// New connects a client to an NNTP server.
func New(net, addr string) (*Client, error) {
	conn, err := textproto.Dial(net, addr)
	if err != nil {
		return nil, err
	}

	return connect(conn)
}

func connect(conn *textproto.Conn) (*Client, error) {
	_, msg, err := conn.ReadCodeLine(200)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		Banner: msg,
	}, nil
}

// Close this client.
func (c *Client) Close() error {
	return c.conn.Close()
}

func parsePosting(p string) nntp.PostingStatus {
	switch p {
	case "y":
		return nntp.PostingPermitted
	case "m":
		return nntp.PostingModerated
	}
	return nntp.PostingNotPermitted
}

// List returns a list of possible groups
func (c *Client) List() (io.Reader, error) {
	err := c.conn.PrintfLine("LIST")
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(215)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Group selects a group.
func (c *Client) Group(name string) (*nntp.Group, error) {

	err := c.conn.PrintfLine("GROUP " + name)
	if err != nil {
		return nil, err
	}
	_, msg, err := c.conn.ReadCodeLine(211)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(msg, " ")
	if len(parts) != 4 {
		err = errors.New("Don't know how to parse result: " + msg)
	}

	var rv nntp.Group
	rv.Count, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	rv.Low, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	rv.High, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, err
	}
	rv.Name = parts[3]

	return &rv, nil
}

// Article grabs an article
func (c *Client) Article(specifier string) (io.Reader, error) {
	err := c.conn.PrintfLine("ARTICLE %s", specifier)
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(220)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Head gets the headers for an article
func (c *Client) Head(specifier string) (io.Reader, error) {
	err := c.conn.PrintfLine("HEAD %s", specifier)
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(221)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Help returns a short summary of commands that are understood by implementation of current server
func (c *Client) Help() (io.Reader, error) {
	err := c.conn.PrintfLine("HELP")
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(100)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Newgroups returns a list of groups createn from given date+time
func (c *Client) Newgroups(date string, time string) (io.Reader, error) {
	err := c.conn.PrintfLine("NEWGROUPS %s %s", date, time)
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(231)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Newnews returns all new news from dat to time in a specified groups
func (c *Client) Newnews(newsgroups string, date string, time string) (io.Reader, error) {
	err := c.conn.PrintfLine("NEWNEWS %s %s %s", newsgroups, date, time)
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(231)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Body gets the body of an article
func (c *Client) Body(specifier string) (io.Reader, error) {
	err := c.conn.PrintfLine("BODY %s", specifier)
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(222)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Stat gets the statistics of an article
func (c *Client) Stat(specifier string) (*string, error) {
	err := c.conn.PrintfLine("STAT %s", specifier)
	if err != nil {
		return nil, err
	}
	_, msg, err := c.conn.ReadCodeLine(223)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Last the internally maintained "current article pointer" is set to the previous article in the current newsgroup.
func (c *Client) Last() (io.Reader, error) {
	err := c.conn.PrintfLine("LAST")
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(223)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Next the internally maintained "current article pointer" is set to the next article in the current newsgroup.
func (c *Client) Next() (io.Reader, error) {
	err := c.conn.PrintfLine("NEXT")
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(233)
	if err != nil {
		return nil, err
	}
	return c.conn.DotReader(), nil
}

// Quit finsishes session with the server
func (c *Client) Quit() (*string, error) {
	err := c.conn.PrintfLine("QUIT")
	if err != nil {
		return nil, err
	}
	_, msg, err := c.conn.ReadCodeLine(205)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Post a new article
func (c *Client) Post(r io.Reader) (*string, error) {
	err := c.conn.PrintfLine("POST")
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(340)
	if err != nil {
		return nil, err
	}
	w := c.conn.DotWriter()
	_, err = io.Copy(w, r)
	if err != nil {
		return nil, err
	}
	w.Close()
	_, msg, err := c.conn.ReadCodeLine(240)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Ihave command IHAVE implementation
func (c *Client) Ihave(r io.Reader, specifier string) (*string, error) {
	err := c.conn.PrintfLine("IHAVE %s", specifier)
	if err != nil {
		return nil, err
	}
	_, _, err = c.conn.ReadCodeLine(335)

	if err != nil {
		return nil, err
	}
	w := c.conn.DotWriter()
	_, err = io.Copy(w, r)
	if err != nil {
		return nil, err
	}
	w.Close()

	_, msg, err := c.conn.ReadCodeLine(235)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
