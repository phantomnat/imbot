package bloon_td6

import (
	"image"

	"gocv.io/x/gocv"
)

type ROIMonkeyPanel struct {
	InfoIcon image.Rectangle

	PtCloseBth     image.Point
	CamoFocusIcon  image.Rectangle
	PtCamoFocusBtn image.Point
	PtPrevFocusBtn image.Point
	PtNextFocusBtn image.Point

	UpgradePath1 ROIMonkeyUpgradePath
	UpgradePath2 ROIMonkeyUpgradePath
	UpgradePath3 ROIMonkeyUpgradePath

	PtSellBtn image.Point
}

type ROIMonkeyUpgradePath struct {
	Levels       [5]image.Rectangle
	Buyable      image.Rectangle
	PtBtnUpgrade image.Point
}

var (
	roiSettingIcon  = Rect(1043, 4, 44, 44)
	roiSettingIcon2 = Rect(780, 4, 44, 44)

	ptsShopDragDown = [2]image.Point{image.Pt(1180, 620), image.Pt(1180, 5)}
	ptsShopDragUp   = [2]image.Point{image.Pt(1180, 105), image.Pt(1180, 715)}

	ptBtnPlay          = image.Pt(1220, 680)
	roiPlaySpeed       = Rect(1222, 675, 2, 4)
	roiMonkeyShopPanel = Rect(1100, 100, 156, 533)

	thUpgradeBtn = gocv.NewScalar(0, 200, 100, 0)

	roiLeftMonkeyPanel = ROIMonkeyPanel{
		InfoIcon:      Rect(38, 95, 37, 37),
		CamoFocusIcon: Rect(213, 166, 50, 49),

		UpgradePath1: ROIMonkeyUpgradePath{
			Levels: [5]image.Rectangle{
				Rect(30, 353, 15, 15),
				Rect(30, 336, 15, 15),
				Rect(30, 320, 15, 15),
				Rect(30, 304, 15, 15),
				Rect(30, 287, 15, 15),
			},
			Buyable:      Rect(170, 323, 5, 5),
			PtBtnUpgrade: image.Pt(220, 325),
		},
		UpgradePath2: ROIMonkeyUpgradePath{
			Levels: [5]image.Rectangle{
				Rect(30, 453, 15, 15),
				Rect(30, 436, 15, 15),
				Rect(30, 420, 15, 15),
				Rect(30, 404, 15, 15),
				Rect(30, 387, 15, 15),
			},
			Buyable:      Rect(170, 423, 5, 5),
			PtBtnUpgrade: image.Pt(220, 425),
		},
		UpgradePath3: ROIMonkeyUpgradePath{
			Levels: [5]image.Rectangle{
				Rect(30, 552, 15, 15),
				Rect(30, 536, 15, 15),
				Rect(30, 520, 15, 15),
				Rect(30, 503, 15, 15),
				Rect(30, 487, 15, 15),
			},
			Buyable:      Rect(170, 523, 5, 5),
			PtBtnUpgrade: image.Pt(220, 525),
		},

		PtCloseBth:     image.Pt(265, 50),
		PtCamoFocusBtn: image.Pt(240, 190),
		PtPrevFocusBtn: image.Pt(55, 250),
		PtNextFocusBtn: image.Pt(240, 250),
		PtSellBtn:      image.Pt(220, 600),
	}

	roiRightMonkeyPanel = ROIMonkeyPanel{
		InfoIcon:      Rect(851, 94, 37, 37),
		CamoFocusIcon: Rect(1027, 166, 50, 49),

		UpgradePath1: ROIMonkeyUpgradePath{
			Levels: [5]image.Rectangle{
				Rect(844, 353, 15, 15),
				Rect(844, 336, 15, 15),
				Rect(844, 320, 15, 15),
				Rect(844, 304, 15, 15),
				Rect(844, 287, 15, 15),
			},
			Buyable:      Rect(985, 323, 5, 5),
			PtBtnUpgrade: image.Pt(1030, 325),
		},
		UpgradePath2: ROIMonkeyUpgradePath{
			Levels: [5]image.Rectangle{
				Rect(844, 453, 15, 15),
				Rect(844, 436, 15, 15),
				Rect(844, 420, 15, 15),
				Rect(844, 404, 15, 15),
				Rect(844, 387, 15, 15),
			},
			Buyable:      Rect(985, 423, 5, 5),
			PtBtnUpgrade: image.Pt(1030, 425),
		},
		UpgradePath3: ROIMonkeyUpgradePath{
			Levels: [5]image.Rectangle{
				Rect(844, 552, 15, 15),
				Rect(844, 536, 15, 15),
				Rect(844, 520, 15, 15),
				Rect(844, 503, 15, 15),
				Rect(844, 487, 15, 15),
			},
			Buyable:      Rect(985, 523, 5, 5),
			PtBtnUpgrade: image.Pt(1030, 525),
		},

		PtCloseBth:     image.Pt(1080, 50),
		PtCamoFocusBtn: image.Pt(1050, 190),
		PtPrevFocusBtn: image.Pt(870, 250),
		PtNextFocusBtn: image.Pt(1050, 250),
		PtSellBtn:      image.Pt(1030, 600),
	}
)
