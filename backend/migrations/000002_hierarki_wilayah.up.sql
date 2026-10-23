CREATE FUNCTION hierarki_wilayah(kode_wilayah text) RETURNS json
LANGUAGE sql STABLE AS $$
    SELECT COALESCE(
        json_agg(json_build_object('kode', kode, 'nama', nama, 'tingkat', tingkat) ORDER BY length(kode)),
        '[]'::json
    )
    FROM wilayah
    WHERE kode_wilayah LIKE kode || '.%'
$$;
