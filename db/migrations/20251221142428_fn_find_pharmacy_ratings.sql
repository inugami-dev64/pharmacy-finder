-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION find_pharmacy_ratings(_id BIGINT)
RETURNS TABLE (
	id BIGINT,
	stars REAL,
	hrt_kind BPCHAR
) AS $$
	SELECT
		pr.pharmacy_id AS id,
		AVG(pr.stars) AS stars,
		NULL AS hrt_kind
	FROM
		pharmacy_reviews pr
	WHERE
		pr.pharmacy_id = _id
	GROUP BY
		pr.pharmacy_id
	UNION ALL
	SELECT
		pr.pharmacy_id AS id,
		AVG(pr.stars) AS stars,
		pr.hrt_kind
	FROM
		pharmacy_reviews pr
	WHERE
		pr.pharmacy_id = _id
	AND
		pr.hrt_kind = 'e'
	GROUP BY
		pr.pharmacy_id,
		pr.hrt_kind
	UNION ALL
	SELECT
		pr.pharmacy_id AS id,
		AVG(pr.stars) AS stars,
		pr.hrt_kind
	FROM
		pharmacy_reviews pr
	WHERE
		pr.pharmacy_id = _id
	AND
		pr.hrt_kind = 't'
	GROUP BY
		pr.pharmacy_id,
		pr.hrt_kind
$$ LANGUAGE SQL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION find_pharmacy_ratings;
-- +goose StatementEnd
