/*
 * Serve /robots.txt
 */
package handler

const text = `User-agent: *
Allow: /`

type HandlerRobots struct{}
