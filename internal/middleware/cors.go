package middleware

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CORS(allowSite string) fiber.Handler {
	allowSite = strings.TrimSpace(allowSite)

	// Jika ALLOW_SITE kosong:
	// hanya izinkan localhost / 127.0.0.1 semua port.
	if allowSite == "" {
		return cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				return isLocalhost(origin)
			},
			AllowMethods: []string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodPut,
				fiber.MethodDelete,
				fiber.MethodPatch,
				fiber.MethodHead,
				fiber.MethodOptions,
			},
			AllowHeaders: []string{
				"Origin",
				"Content-Type",
				"Accept",
				"Authorization",
			},
		})
	}

	allowedSites := splitOrigins(allowSite)

	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return isAllowedOrigin(
				origin,
				allowedSites,
			)
		},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodHead,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	})
}

func isLocalhost(origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if parsed.Scheme != "http" &&
		parsed.Scheme != "https" {
		return false
	}

	host := strings.ToLower(
		parsed.Hostname(),
	)

	return host == "localhost" ||
		host == "127.0.0.1"
}

func isAllowedOrigin(
	origin string,
	allowedSites []string,
) bool {

	origin = strings.TrimSpace(origin)

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	originHost := strings.ToLower(
		parsed.Hostname(),
	)

	originScheme := strings.ToLower(
		parsed.Scheme,
	)

	for _, allowed := range allowedSites {
		allowed = strings.TrimSpace(
			allowed,
		)

		if allowed == "" {
			continue
		}

		// Wildcard subdomain:
		//
		// https://*.micinproject.my.id
		//
		if strings.HasPrefix(
			allowed,
			"https://*.",
		) || strings.HasPrefix(
			allowed,
			"http://*.",
		) {

			allowedURL, err := url.Parse(
				allowed,
			)
			if err != nil {
				continue
			}

			allowedHost := strings.TrimPrefix(
				strings.ToLower(
					allowedURL.Hostname(),
				),
				"*.",
			)

			if originScheme != strings.ToLower(
				allowedURL.Scheme,
			) {
				continue
			}

			// Harus benar-benar subdomain.
			//
			// app.micinproject.my.id  ✅
			// api.micinproject.my.id  ✅
			// micinproject.my.id      ❌
			if strings.HasSuffix(
				originHost,
				"."+allowedHost,
			) {
				return true
			}

			continue
		}

		// Exact origin.
		if strings.EqualFold(
			origin,
			allowed,
		) {
			return true
		}
	}

	return false
}

func splitOrigins(value string) []string {
	parts := strings.Split(
		value,
		",",
	)

	origins := make(
		[]string,
		0,
		len(parts),
	)

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		origins = append(
			origins,
			part,
		)
	}

	return origins
}
